# 背景
这是一个密码管理器的web service后端代码

提供的所有接口：

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

# 你的任务
- 创建一个基于vue3的前端项目，自行设计页面结构，对接上面所有的API。页面风格清爽简洁，并且必须符合用户的使用直觉
- 在请求`/pwdManage/saveOne`之前，需要弹窗提示用户输入一个 字符串 作为对称加密的密钥，**将用户待保存的密码进行加密后的结果**作为Password字段的值传递到后端持久化
- 在获取到`/pwdManage/query`的返回结果之后，密码栏脱敏显示为`***`，提供一个按钮，点击后输入解密密码，才展示给用户
- 可以在本地启动这个前端项目并调试
- 构造docker file，直接容器化打包部署