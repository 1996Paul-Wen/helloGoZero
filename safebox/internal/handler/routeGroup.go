package handler

import (
	"net/http"

	"github.com/1996Paul-Wen/helloGoZero/safebox/internal/svc"
	"github.com/zeromicro/go-zero/rest"
)

type RouteGroup struct {
	Prefix      string
	Routes      []rest.Route
	svcCtx      *svc.ServiceContext
	middlewares []rest.Middleware
}

func NewRouteGroup(prefix string, svcCtx *svc.ServiceContext) *RouteGroup {
	return &RouteGroup{
		Prefix: prefix,
		svcCtx: svcCtx,
	}
}

func (rg *RouteGroup) POST(path string, handler http.HandlerFunc) {
	rg.Routes = append(rg.Routes, rest.Route{
		Method:  http.MethodPost,
		Path:    path,
		Handler: handler,
	})
}

func (rg *RouteGroup) GET(path string, handler http.HandlerFunc) {
	rg.Routes = append(rg.Routes, rest.Route{
		Method:  http.MethodGet,
		Path:    path,
		Handler: handler,
	})
}

func (rg *RouteGroup) AddMiddleware(middleWare rest.Middleware) {
	rg.middlewares = append(rg.middlewares, middleWare)
}

func (rg *RouteGroup) RegisterToServer(server *rest.Server) {
	withMiddleWareRoutes := rest.WithMiddlewares(rg.middlewares, rg.Routes...)
	server.AddRoutes(withMiddleWareRoutes, rest.WithPrefix(rg.Prefix))
}
