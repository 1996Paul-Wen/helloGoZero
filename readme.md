# safebox基于go-zero初始化
本项目基于go-zero框架初始化：
```bash
# install goctl
go install github.com/zeromicro/go-zero/tools/goctl@latest

# init go module
go mod init github.com/1996Paul-Wen/helloGoZero
# create framework by goctl
goctl api new greet
# install dependencies
go mod tity


# run server
go run ./greet/greet.go -f ./greet/etc/greet-api.yaml

# test service 
curl -i http://localhost:8888/greet/from/you
```

基于该框架，搭建safebox

safebox应用启动命令：`go run ./safebox/safebox.go -f ./safebox/etc/safebox-api.yaml`


# 项目依赖的组件
- MariaDB: `10.3.39-MariaDB-deepin1`. 可更换为mysql
- ES: 7.17.29

# 建表ddl
```sql
CREATE DATABASE IF NOT EXISTS `safebox` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

use `safebox`;

CREATE TABLE `user` (
    `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `username`    VARCHAR(50) NOT NULL COMMENT '用户名',
    `hash_password` VARCHAR(255) NOT NULL COMMENT 'bcrypt哈希后的密码', 
    `creator`     VARCHAR(50) NOT NULL COMMENT '创建人',
    `updator`     VARCHAR(50) NOT NULL COMMENT '更新人',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';


CREATE TABLE `managed_password` (
    `id`          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `user_id`     BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `description` VARCHAR(255) NOT NULL COMMENT '对于本条账号密码的明文描述，将用于es倒排',
    `username`    VARCHAR(50) NOT NULL COMMENT '用户名（对称加密后的值）',
    `password` VARCHAR(255) NOT NULL COMMENT '密码（对称加密后的值）', 
    `creator`     VARCHAR(50) NOT NULL COMMENT '创建人',
    `updator`     VARCHAR(50) NOT NULL COMMENT '更新人',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    INDEX `user_id_idx` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='被管理的用户密码表';
```


# 生成model模板代码
```bash
goctl model mysql datasource -url="root:abcd@tcp(127.0.0.1:3306)/safebox" -table="user" -dir="./safebox/internal/model" -cache=false --style=goZero
goctl model mysql datasource -url="root:abcd@tcp(127.0.0.1:3306)/safebox" -table="managed_password" -dir="./safebox/internal/model" -cache=false --style=goZero
```

# curl
```bash
# create user
curl -POST http://localhost:8888/user/create -d '{"name": "efg", "password": "456"}' -H "Content-Type: application/json"

resp:
{
    "Code": 0,
    "Msg": "",
    "Data": {
        "UserID": 1
    },
    "TraceID": "a5bf75e7-dc32-4c89-a38f-c125c7606777"
}
```

```bash

# login
curl -POST http://localhost:8888/user/login -d '{"name": "efg", "password": "456"}' -H "Content-Type: application/json"

resp:
{
    "Code": 0,
    "Msg": "",
    "Data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzU5OTAyMDIsImlhdCI6MTc3NTk4OTkwMiwidXNlcklkIjoxfQ.H91Tp6_cI4WnSvWIbMYnnGPHCBYWHSaN_EIzmnSeXOE"
    },
    "TraceID": "3a4e8a33-c997-4f8f-bb0b-798556e2c466"
}
```

```bash

# describe user
curl -POST http://localhost:8888/user/describe -d '{}' -H "Content-Type: application/json" -H "Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzU5OTAyMDIsImlhdCI6MTc3NTk4OTkwMiwidXNlcklkIjoxfQ.H91Tp6_cI4WnSvWIbMYnnGPHCBYWHSaN_EIzmnSeXOE"

resp:
{
    "Code": 0,
    "Msg": "",
    "Data": {
        "Id": 1,
        "Username": "efg",
        "HashPassword": "***",
        "Creator": "efg",
        "Updator": "efg",
        "CreateTime": "2026-04-12T10:30:23+08:00",
        "UpdateTime": "2026-04-12T10:30:23+08:00"
    },
    "TraceID": "7e167b92-5712-4190-8511-8b1d6f462e1b"
}
```

```bash
# 插入一条用户名和密码
curl -POST http://localhost:8888/pwdManage/saveOne -d '{"description":"测试用中信银行非上海银行手机app","username":"testabc", "password": "testpwd"}' -H "Content-Type: application/json" -H "Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzU5OTAyMDIsImlhdCI6MTc3NTk4OTkwMiwidXNlcklkIjoxfQ.H91Tp6_cI4WnSvWIbMYnnGPHCBYWHSaN_EIzmnSeXOE"

resp:
{
    "Code": 0,
    "Msg": "",
    "Data": null,
    "TraceID": "5be786dd-1a77-4fc5-8d31-44cf5e0e2a9b"
}
```

```bash
# 查询可能的用户名和及其密码
curl -POST http://localhost:8888/pwdManage/query -d '{"query":"中信 银行"}' -H "Content-Type: application/json" -H "Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzU5OTAyMDIsImlhdCI6MTc3NTk4OTkwMiwidXNlcklkIjoxfQ.H91Tp6_cI4WnSvWIbMYnnGPHCBYWHSaN_EIzmnSeXOE"

resp:
{
    "Code": 0,
    "Msg": "",
    "Data": [
        {
            "Id": 1,
            "UserId": 1,
            "Description": "测试用中信银行非上海银行手机app",
            "Username": "testabc",
            "Password": "testpwd",
            "Creator": "efg",
            "Updator": "efg",
            "CreateTime": "2026-04-12T10:34:13+08:00",
            "UpdateTime": "2026-04-12T10:34:13+08:00"
        }
    ],
    "TraceID": "40a5b1f7-0eea-4983-80d6-36be08250e8b"
}
```

```bash
# 登陆态过期查询

```bash
 ~ % curl -POST http://127.0.0.1:8888/user/describe -d '{}' -H "Content-Type: application/json" -H "Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzM2NDU2OTksImlhdCI6MTc3MzY0NTM5OSwidXNlcklkIjoxfQ.UALzZWS1_OBYlpyOnOoxJZICOM5bEqs8NaRWbaL65Yk" -v
*   Trying 127.0.0.1:8888...
* Connected to 127.0.0.1 (127.0.0.1) port 8888
> POST /user/describe HTTP/1.1
> Host: 127.0.0.1:8888
> User-Agent: curl/8.7.1
> Accept: */*
> Content-Type: application/json
> Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzM2NDU2OTksImlhdCI6MTc3MzY0NTM5OSwidXNlcklkIjoxfQ.UALzZWS1_OBYlpyOnOoxJZICOM5bEqs8NaRWbaL65Yk
> Content-Length: 2
> 
* upload completely sent off: 2 bytes
< HTTP/1.1 401 Unauthorized
< Content-Type: text/plain; charset=utf-8
< Traceparent: 00-aaf035959e9e46d4c18fafced17c7c90-c200a12fa6c057b4-00
< X-Content-Type-Options: nosniff
< Date: Sun, 12 Apr 2026 10:28:51 GMT
< Content-Length: 13
< 
Unauthorized
* Connection #0 to host 127.0.0.1 left intact
 ~ % 

```