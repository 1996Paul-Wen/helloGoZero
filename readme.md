# safebox后端基于go-zero初始化
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
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
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
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
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
# 修改已保存的用户名/密码
curl -POST http://localhost:8888/pwdManage/updateOne -d '{"id":1, "description":"测试用中信银行非上海银行手机app", "username":"testabc", "password": "testpwd"}' -H "Content-Type: application/json" -H "Content-Type: application/json" -H "Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzU5OTAyMDIsImlhdCI6MTc3NTk4OTkwMiwidXNlcklkIjoxfQ.H91Tp6_cI4WnSvWIbMYnnGPHCBYWHSaN_EIzmnSeXOE"

resp:
{
    "Code": 0,
    "Msg": "",
    "Data": 1,
    "TraceID": "5be786dd-1a77-4fc5-8d31-44cf5e0e2a9b"
}
```

```bash
# 删除一条用户名和密码
curl -POST http://localhost:8888/pwdManage/deleteOne -d '{"id":1}' -H "Content-Type: application/json" -H "Content-Type: application/json" -H "Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzU5OTAyMDIsImlhdCI6MTc3NTk4OTkwMiwidXNlcklkIjoxfQ.H91Tp6_cI4WnSvWIbMYnnGPHCBYWHSaN_EIzmnSeXOE"

resp:
{
    "Code": 0,
    "Msg": "",
    "Data": null,
    "TraceID": "5be786dd-1a77-4fc5-8d31-44cf5e0e2a9b"
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


# 从固定密钥衍生随机密钥，防止数据库爆破
safebox前端对用户输入的待保存的密码，使用AES对称加密后，将得到的密文密码交给后端入库

为了防止数据库爆破，用户保存的每个密码需要使用不同的AES密钥

但用户不可能每次使用不同的AES密钥去加密密码，否则用户又将需要记忆大量的AES密钥，从而陷入用一个密码加密另一个密码的困境

safebox的解决方案是，
- **一个用户只需要使用一个AES密钥，由safebox通过计算来衍生出不同的、最终的AES密钥**
- **密文密码将插入固定长度的随机前缀，使得相同明文每次加密结果不同，防止模式分析**


## AES 静态密钥 → 动态密钥的渲染机制

核心在 `safebox-web/src/utils/crypto.ts`，整个流程如下：

### 整体架构

```
┌──────────────────────────────────────────────────────────────┐
│                     用户输入（"静态"）                         │
│  passphrase: "我的加密密码123"   ← 用户每次手动输入的主密钥     │
│                                                              │
│                        ↓ 派生算法                            │
│                                                              │
│              deriveKey() 函数                                │
│    ┌─────────────────────────────────────┐                  │
│    │ passphrase + hashSuffix(desc+user)  │                  │
│    │         ↓                           │                  │
│    │ PBKDF2(10000次迭代, 固定Salt)       │                  │
│    │         ↓                           │                  │
│    │ 输出：256-bit AES 密钥（动态）       │                  │
│    └─────────────────────────────────────┘                  │
│                        ↓                                    │
│              AES-256-CBC 加密/解密                          │
└──────────────────────────────────────────────────────────────┘
```

### 核心代码逐层解析

**第一步：生成每条记录唯一的 hash 后缀**

```10:14:safebox-web/src/utils/crypto.ts
function hashSuffix(description: string, username: string): string {
  const combined = description + username          // "Gmail账号"+"zhangsan"
  const hash = CryptoJS.SHA256(combined)            // SHA-256 哈希
  return hash.toString(CryptoJS.enc.Hex).slice(-4)  // 取最后 4 个 hex 字符，如 "a3f2"
}
```

> 这一步让 **不同密码记录产生不同的后缀**，保证每条记录的实际加密密钥不同。

**第二步：PBKDF2 密钥派生（静态 → 动态）**

```20:26:safebox-web/src/utils/crypto.ts
function deriveKey(passphrase: string, description: string, username: string): CryptoJS.lib.WordArray {
  const finalKey = passphrase + hashSuffix(description, username)
  // 例: "我的密码123" + "a3f2" = "我的密码123a3f2"
  
  return CryptoJS.PBKDF2(finalKey, SALT, {        // ← 核心！PBKDF2 拉伸
    keySize: KEY_SIZE,      // 256 bits (32 bytes) = AES-256
    iterations: 10000,      // 10000 次哈希迭代，增加暴力破解成本
  })
}
```

> `SALT = 'SafeBox2024SaltKey'` 是硬编码的固定盐值。

**第三步：AES-256-CBC 加密**

```36:54:safebox-web/src/utils/crypto.ts
export function encrypt(plainText, passphrase, description, username): string {
  const key = deriveKey(passphrase, description, username)  // 动态密钥
  const iv = CryptoJS.lib.WordArray.random(128 / 8)          // 随机 IV (16字节)
  
  const encrypted = CryptoJS.AES.encrypt(plainText, key, {
    iv: iv,
    mode: CryptoJS.mode.CBC,
    padding: CryptoJS.pad.Pkcs7,
  })
  
  // IV + 密文 拼在一起 Base64 编码输出
  const combined = iv.concat(encrypted.ciphertext)
  return combined.toString(CryptoJS.enc.Base64)
}
```

### 具体示例

假设用户保存一条 **Gmail 账号密码**：

| 参数 | 值 |
|------|-----|
| 描述 | `Gmail账号` |
| 用户名 | `zhangsan@gmail.com` |
| 明文密码 | `MyGmailPass2024!` |
| 用户输入的加密密钥 | `mypass123` |

**派生过程**：

```
1. hashSuffix("Gmail账号", "zhangsan@gmail.com")
   → SHA256("Gmail账号zhangsan@gmail.com") 的最后4位 → 例如 "7d3e"

2. finalKey = "mypass123" + "7d3e" = "mypass1237d3e"

3. PBKDF2("mypass1237d3e", "SafeBox2024SaltKey", 10000次)
   → 32字节的 AES-256 动态密钥

4. 用该动态密钥 + 随机IV 做 AES-256-CBC 加密
   → 输出 Base64 密文存入数据库
```

如果保存另一条 **GitHub 账号**（用同一个主密钥）：

```
1. hashSuffix("GitHub账号", "zhangsan") 
   → SHA256 不同 → 后缀不同，例如 "a1b2"

2. finalKey = "mypass123" + "a1b2" = "mypass123a1b2"

3. PBKDF2 输出的密钥完全不同！

→ 同一个主密码 "mypass123"，两条记录的实际 AES 密钥不同 ✓
```

---

## 设计总结

| 层级 | 内容 | 是否存储 |
|------|------|---------|
| **静态密钥** | 用户每次手动输入的 `passphrase`（如 "mypass123"） | ❌ 不存储 |
| **固定盐值** | `SafeBox2024SaltKey`（硬编码在前端代码中） | ✅ 在前端源码中 |
| **记录指纹** | `hashSuffix(description+username)` 后 4 位 | ❌ 不存储，可复算 |
| **动态密钥** | PBKDF2 派生出的 256-bit AES 密钥 | ❌ 不存储，用完即弃 |
| **随机 IV** | 每次加密生成的 16 字节随机向量 | ✅ 存在密文头部 |
| **最终密文** | `Base64(IV + ciphertext)` | ✅ 存入 MySQL |

**关键设计思想**：
1. **一主密钥 + 多条记录 = 多个不同的 AES 密钥** — 即使一条被破解，其他条目不受影响
2. **PBKDF2 10000 次迭代** — 抵御暴力破解/彩虹表攻击
3. **随机 IV** — 相同明文每次加密结果不同，防止模式分析
4. **服务端永远不接触明文和密钥** — 纯前端加密，后端只存密文