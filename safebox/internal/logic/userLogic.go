package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/1996Paul-Wen/helloGoZero/safebox/internal/infra"
	"github.com/1996Paul-Wen/helloGoZero/safebox/internal/model"
	"github.com/1996Paul-Wen/helloGoZero/safebox/internal/svc"
	"github.com/1996Paul-Wen/helloGoZero/safebox/internal/types"
	"github.com/1996Paul-Wen/helloGoZero/safebox/internal/util"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type UserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLogic {
	return &UserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (u *UserLogic) Create(req *types.CreateUserReq) (userID int64, err error) {
	conn := infra.LoadSQLConn()

	userModel := model.NewUserModel(conn)
	dbUser, err := userModel.FindOneByUsername(u.ctx, req.Name)
	if err != nil && err != model.ErrNotFound {
		return
	}
	if dbUser != nil {
		err = fmt.Errorf("%+v has registed", dbUser.Username)
		return
	}

	// bcrypt 加密 password
	// 生成哈希，使用默认成本（10）
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	sqlResult, err := userModel.Insert(u.ctx, &model.User{
		Username:     req.Name,
		HashPassword: string(hashedPassword),
		Creator:      req.Name,
		Updator:      req.Name,
	})
	if err != nil {
		return
	}
	userID, err = sqlResult.LastInsertId()
	return

}

type LoginResp struct {
	Token string `json:"token"`
}

func (u *UserLogic) Login(req *types.UserLoginReq) (loginResp LoginResp, err error) {
	conn := infra.LoadSQLConn()
	userModel := model.NewUserModel(conn)

	dbUser, err := userModel.FindOneByUsername(u.ctx, req.Name)
	if err != nil {
		return
	}
	hashPassword := dbUser.HashPassword

	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(req.Password))
	if err != nil {
		err = fmt.Errorf("密码错误")
		return
	}

	// gen jwt token
	now := time.Now().Unix()
	claims := make(jwt.MapClaims)
	claims[util.JWTEXP] = now + u.svcCtx.Config.Auth.AccessExpire
	claims[util.JWTIAT] = now
	claims[util.JWTKeyUserID] = dbUser.Id

	token, err := util.GenerateToken(u.svcCtx.Config.Auth.AccessSecret, claims)
	if err != nil {
		return
	}
	loginResp.Token = token
	return
}

func (u *UserLogic) Describe() (dbUser *model.User, err error) {
	conn := infra.LoadSQLConn()
	userModel := model.NewUserModel(conn)

	userId := u.ctx.Value(util.JWTKeyUserID)
	dbUser, err = userModel.FindOne(u.ctx, userId.(uint64))
	if err != nil {
		return
	}
	dbUser.HashPassword = "***"
	return
}
