// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf // 内嵌 go-zero 默认的 RESTful 服务配置

	// 自定义配置
	MySQL struct { // 数据库配置
		DataSource string
	}
	ESConfig ESConfigModel
}

type ESConfigModel struct { // ES
	Hosts []string // es 节点地址列表
	// Username string // 可选，基本认证用户名
	// Password string // 可选，基本认证密码
}
