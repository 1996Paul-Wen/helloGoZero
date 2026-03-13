package handler

import (
	"net/http"

	"github.com/1996Paul-Wen/helloGoZero/safebox/internal/logic"
	"github.com/1996Paul-Wen/helloGoZero/safebox/internal/svc"
	"github.com/1996Paul-Wen/helloGoZero/safebox/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var userRouteGroup *RouteGroup

func InitUserRouteGroup(svcCtx *svc.ServiceContext) {
	userRouteGroup = NewRouteGroup("/user", svcCtx)

	jwtAuthMiddleWare := BuildAuthMiddleware(svcCtx)

	userRouteGroup.POST("/create", CreateUser(svcCtx))
	userRouteGroup.POST("/login", Login(svcCtx))
	userRouteGroup.POST("/describe", jwtAuthMiddleWare(Describe(svcCtx)))

}

func CreateUser(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateUserReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		u := logic.NewUserLogic(r.Context(), svcCtx)
		resp, err := u.Create(&req)
		if err != nil {
			// httpx.ErrorCtx(r.Context(), w, err)
			httpx.OkJsonCtx(r.Context(), w, BuildFailResp(r.Context(), -1, err))
		} else {
			httpx.OkJsonCtx(r.Context(), w, BuildSuccessResp(r.Context(), struct {
				UserID int64
			}{
				UserID: resp,
			}))
		}
	}
}

func Login(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserLoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		u := logic.NewUserLogic(r.Context(), svcCtx)
		resp, err := u.Login(&req)
		if err != nil {
			// httpx.ErrorCtx(r.Context(), w, err)
			httpx.OkJsonCtx(r.Context(), w, BuildFailResp(r.Context(), -1, err))
		} else {
			httpx.OkJsonCtx(r.Context(), w, BuildSuccessResp(r.Context(), resp))
		}
	}
}

func Describe(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := logic.NewUserLogic(r.Context(), svcCtx)
		resp, err := u.Describe()
		if err != nil {
			// httpx.ErrorCtx(r.Context(), w, err)
			httpx.OkJsonCtx(r.Context(), w, BuildFailResp(r.Context(), -1, err))
		} else {
			httpx.OkJsonCtx(r.Context(), w, BuildSuccessResp(r.Context(), resp))
		}
	}
}
