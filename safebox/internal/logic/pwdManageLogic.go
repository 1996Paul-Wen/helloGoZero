package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/1996Paul-Wen/helloGoZero/safebox/internal/infra"
	"github.com/1996Paul-Wen/helloGoZero/safebox/internal/model"
	"github.com/1996Paul-Wen/helloGoZero/safebox/internal/svc"
	"github.com/1996Paul-Wen/helloGoZero/safebox/internal/util"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var managedPWDESIndexName = "managed_passwords"

type PWDManageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPWDManageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PWDManageLogic {
	return &PWDManageLogic{Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx}
}

type SavePWD struct {
	Description string `json:"description" validate:"required"`
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password" validate:"required"`
}

func (pml *PWDManageLogic) SaveOne(savePWD SavePWD) (int64, error) {
	// validate savePWD
	validate := validator.New()
	err := validate.Struct(savePWD)
	if err != nil {
		// 处理验证错误
		var errString strings.Builder
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Fprintf(&errString, "字段 %s 不能为空\n", err.Field())
		}
		return -1, errors.New(errString.String())
	}

	userId := pml.ctx.Value(util.JWTKeyUserID)
	userIdUint, ok := userId.(uint64)
	if !ok {
		return -1, errors.New("invalid user id")
	}

	conn := infra.LoadSQLConn()
	var insertedId int64
	// 启动事务 编程式事务
	err = conn.Transact(func(session sqlx.Session) error {
		// 在事务内使用 session 初始化模型，确保共用同一事务
		userModel := model.NewUserModel(sqlx.NewSqlConnFromSession(session))
		pwdManageModel := model.NewManagedPasswordModel(sqlx.NewSqlConnFromSession(session))

		userDetail, err := userModel.FindOne(pml.ctx, userIdUint)
		if err != nil {
			return err
		}

		dataToInsert := &model.ManagedPassword{
			UserId:      userIdUint,
			Description: savePWD.Description,
			Username:    savePWD.Username,
			Password:    savePWD.Password,
			Creator:     userDetail.Username,
			Updator:     userDetail.Username,
		}

		sqlResult, err := pwdManageModel.Insert(pml.ctx, dataToInsert)
		if err != nil {
			return err
		}

		insertedId, err = sqlResult.LastInsertId()
		if err != nil {
			return err
		}

		// ---存入es---
		esClient := infra.LoadESClient()
		// 构造文档内容：description + username，方便模糊搜索
		doc := map[string]interface{}{
			"content": savePWD.Description + " " + savePWD.Username,
		}
		data, err := json.Marshal(doc)
		if err != nil {
			return err
		}

		// 索引文档，使用 MySQL 自增 ID 作为文档 ID
		resp, err := esClient.Index(
			managedPWDESIndexName, // 索引名称
			bytes.NewReader(data),
			esClient.Index.WithDocumentID(fmt.Sprintf("%d", insertedId)),
			esClient.Index.WithContext(pml.ctx),
			esClient.Index.WithRefresh("true"), // 立即刷新，方便后续查询
		)
		if err != nil {
			return fmt.Errorf("ES 索引失败: %w", err)
		}
		defer resp.Body.Close()
		if resp.IsError() {
			return fmt.Errorf("ES 索引失败: %s", resp.String())
		}

		return nil
	})

	return insertedId, err

	// userModel := model.NewUserModel(conn)
	// pwdManageModel := model.NewManagedPasswordModel(conn)

	// userDetail, err := userModel.FindOne(pml.ctx, userIdUint)
	// if err != nil {
	// 	return -1, err
	// }

	// dataToInsert := &model.ManagedPassword{
	// 	UserId:      userId.(uint64),
	// 	Description: savePWD.Description,
	// 	Username:    savePWD.Username,
	// 	Password:    savePWD.Password,
	// 	Creator:     userDetail.Username,
	// 	Updator:     userDetail.Username,
	// }
	// sqlResult, err := pwdManageModel.Insert(pml.ctx, dataToInsert)
	// if err != nil {
	// 	return -1, err
	// }
	// insertedId, err := sqlResult.LastInsertId()
	// return insertedId, err

}

// QueryByES 根据关键词在 Elasticsearch 中搜索，返回匹配的密码管理记录
func (pml *PWDManageLogic) QueryByES(query string) ([]*model.ManagedPassword, error) {
	userId := pml.ctx.Value(util.JWTKeyUserID)
	userIdUint, ok := userId.(uint64)
	if !ok {
		return nil, errors.New("invalid user id")
	}

	// 创建 ES 客户端
	esClient := infra.LoadESClient()

	// 构造 match 查询，搜索 content 字段
	searchBody := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"content": query,
			},
		},
	}
	data, err := json.Marshal(searchBody)
	if err != nil {
		return nil, err
	}

	// 执行搜索 - 修正方式
	resp, err := esClient.Search(
		esClient.Search.WithIndex(managedPWDESIndexName),
		esClient.Search.WithBody(bytes.NewReader(data)),
		esClient.Search.WithContext(pml.ctx),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return nil, fmt.Errorf("ES 搜索失败: %s", resp.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// 解析 hits，获取文档 ID
	hits, ok := result["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok || len(hits) == 0 {
		return nil, nil
	}

	var ids []uint64
	for _, hit := range hits {
		hitMap := hit.(map[string]interface{})
		idStr, ok := hitMap["_id"].(string)
		if !ok {
			continue
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			continue // 忽略无法转换的 ID
		}
		ids = append(ids, uint64(id))
	}

	if len(ids) == 0 {
		return nil, nil
	}

	// 根据 ID 从 MySQL 查询完整记录
	conn := infra.LoadSQLConn()
	pwdManageModel := model.NewManagedPasswordModel(conn)

	var passwords []*model.ManagedPassword
	passwords, err = pwdManageModel.FindByCond(pml.ctx, model.ListPWDCond{
		UserIDs: []uint64{userIdUint},
		IDs:     ids,
	})

	return passwords, err
}
