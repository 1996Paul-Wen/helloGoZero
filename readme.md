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
curl -i http://localhost:8888/from/you
```

基于该框架，搭建safebox

safebox应用启动命令：`go run ./safebox/safebox.go -f ./safebox/etc/safebox-api.yaml`


# 项目依赖的组件
- MariaDB: `10.3.39-MariaDB-deepin1`. 可更换为mysql
- ES: 7.17.29

# 建表ddl
```sql
CREATE DATABASE IF NOT EXISTS `safebox` 
    -> CHARACTER SET utf8mb4 
    -> COLLATE utf8mb4_unicode_ci;

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