// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"flag"
	"fmt"

	"github.com/1996Paul-Wen/helloGoZero/safebox/internal/config"
	"github.com/1996Paul-Wen/helloGoZero/safebox/internal/handler"
	"github.com/1996Paul-Wen/helloGoZero/safebox/internal/infra"
	"github.com/1996Paul-Wen/helloGoZero/safebox/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/safebox-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// init infra
	infra.InitConn(c.MySQL.DataSource)
	infra.InitEsClient(c.ESConfig)

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
