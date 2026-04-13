package handler

import (
	"fmt"
	"net/http"

	"github.com/1996Paul-Wen/helloGoZero/safebox/internal/logic"
	"github.com/1996Paul-Wen/helloGoZero/safebox/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var pwdManageRouteGroup *RouteGroup

func InitPWDManageRouteGroup(svcCtx *svc.ServiceContext) {
	pwdManageRouteGroup = NewRouteGroup("/pwdManage", svcCtx)

	jwtAuthMiddleWare := BuildAuthMiddleware(svcCtx)

	pwdManageRouteGroup.POST("/query", jwtAuthMiddleWare(QueryPWD(svcCtx)))
	pwdManageRouteGroup.POST("/saveOne", jwtAuthMiddleWare(SavePWD(svcCtx)))
	pwdManageRouteGroup.POST("/updateOne", jwtAuthMiddleWare(UpdatePWD(svcCtx)))
	pwdManageRouteGroup.POST("/deleteOne", jwtAuthMiddleWare(DeletePWD(svcCtx)))

}

func SavePWD(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req logic.SavePWD
		if err := httpx.Parse(r, &req); err != nil {
			fmt.Println(err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		p := logic.NewPWDManageLogic(r.Context(), svcCtx)
		_, err := p.SaveOne(req)
		if err != nil {
			// httpx.ErrorCtx(r.Context(), w, err)
			httpx.OkJsonCtx(r.Context(), w, BuildFailResp(r.Context(), -1, err))
		} else {
			httpx.OkJsonCtx(r.Context(), w, BuildSuccessResp(r.Context(), nil))
		}
	}
}

func QueryPWD(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Query string `json:"query"`
		}
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		p := logic.NewPWDManageLogic(r.Context(), svcCtx)
		resp, err := p.QueryByES(req.Query)
		if err != nil {
			// httpx.ErrorCtx(r.Context(), w, err)
			httpx.OkJsonCtx(r.Context(), w, BuildFailResp(r.Context(), -1, err))
		} else {
			httpx.OkJsonCtx(r.Context(), w, BuildSuccessResp(r.Context(), resp))
		}
	}
}

func UpdatePWD(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req logic.SavePWD
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		p := logic.NewPWDManageLogic(r.Context(), svcCtx)
		resp, err := p.UpdateOne(req)
		if err != nil {
			// httpx.ErrorCtx(r.Context(), w, err)
			httpx.OkJsonCtx(r.Context(), w, BuildFailResp(r.Context(), -1, err))
		} else {
			httpx.OkJsonCtx(r.Context(), w, BuildSuccessResp(r.Context(), resp))
		}
	}
}

func DeletePWD(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID int64 `json:"id"`
		}
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		p := logic.NewPWDManageLogic(r.Context(), svcCtx)
		err := p.DeleteOne(req.ID)
		if err != nil {
			// httpx.ErrorCtx(r.Context(), w, err)
			httpx.OkJsonCtx(r.Context(), w, BuildFailResp(r.Context(), -1, err))
		} else {
			httpx.OkJsonCtx(r.Context(), w, BuildSuccessResp(r.Context(), nil))
		}
	}
}
