package infra

import (
	"fmt"
	"log"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var conn sqlx.SqlConn

func LoadSQLConn() sqlx.SqlConn {
	if conn == nil {
		log.Fatal("db conn not init")
	}

	return conn
}

func InitConn(dataSource string) {
	conn = sqlx.NewMysql(dataSource)
	// 执行一个简单的查询，比如检查数据库版本
	var version string
	err := conn.QueryRow(&version, "SELECT VERSION()")
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
	fmt.Printf("Database version: %s\n", version)
}
