// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"github.com/1996Paul-Wen/helloGoZero/safebox/internal/logic"
	"github.com/1996Paul-Wen/helloGoZero/safebox/internal/svc"
	"github.com/1996Paul-Wen/helloGoZero/safebox/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var greetRouteGroup *RouteGroup

func InitGreetRouteGroup(svcCtx *svc.ServiceContext) {
	greetRouteGroup = NewRouteGroup("/greet", svcCtx)

	greetRouteGroup.GET("/from/:name", GreetHandler(svcCtx))
}

func GreetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGreetLogic(r.Context(), svcCtx)
		resp, err := l.Greet(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, BuildSuccessResp(r.Context(), resp))
		}
	}
}
