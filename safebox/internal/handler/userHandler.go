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

	userRouteGroup.POST("/create", CreateUser(svcCtx))

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
