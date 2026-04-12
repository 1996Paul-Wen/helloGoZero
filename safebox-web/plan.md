
Now I have a thorough understanding of the backend API:

1. **User APIs:**
   - `POST /user/create` - Create user with `{name, password}`
   - `POST /user/login` - Login, returns JWT token
   - `POST /user/describe` - Get current user info (requires auth)

2. **Password Management APIs:**
   - `POST /pwdManage/saveOne` - Save a password record `{description, username, password}` (requires auth)
   - `POST /pwdManage/query` - Query passwords by keyword `{query}` (returns encrypted passwords)

3. **Auth:** Bearer token in Authorization header

4. **Frontend requirements:**
   - Vue 3 project
   - Clean and intuitive UI
   - Before saving password: popup to enter symmetric encryption key, encrypt the password before sending
   - After querying: show passwords as `***`, provide button to decrypt after entering key

The plan has been created successfully with 5 main tasks:

1. 初始化 Vue3 + Vite + TypeScript + Tailwind + Element Plus 项目脚手架
2. 搭建路由系统、Pinia authStore、Axios 请求封装及全部 API 接口
3. 实现 AES 加解密工具模块（crypto.ts），含 PBKDF2 密钥派生和 CBC 模式编解码
4. 实现登录页和注册页 UI 及完整交互逻辑
5. 实现仪表盘主页、密码卡片列表、搜索、新增密码（含加密弹窗）、查看密码（含解密弹窗）全套功能

方案概览：

## SafeBox 密码管理器 - Vue3 前端

### 技术选型
- **Vue 3** + **Vite** + **TypeScript**
- **Tailwind CSS** (清爽简洁的样式)
- **Element Plus** (UI 组件库：弹窗、表单、消息提示等)
- **Pinia** (状态管理：token 持久化)
- **Vue Router** (路由)
- **Axios** (HTTP 客户端，带拦截器处理 token 和统一响应格式)
- **crypto-js / Web Crypto API** (AES 对称加解密)

### 页面结构设计
| 页面 | 功能 |
|------|------|
| `/login` | 登录表单 → 获取 token → 跳转主页 |
| `/register` | 注册表单 → 创建用户 → 跳转登录 |
| `/dashboard` | 仪表盘主页 |

### 核心交互流程
1. **保存密码**: 填写描述/用户名/密码 → 点击保存 → 弹出"输入加密密钥"对话框 → AES 加密密码 → 发送到后端
2. **查询密码**: 输入关键词搜索 → 展示卡片列表(密码显示 `***`) → 点击"查看密码"按钮 → 弹出"输入解密密钥"对话框 → 解密后展示明文
3. **Token 管理**: Axios 请求拦截器自动附加 Bearer Token，401 自动跳转登录

### 任务分解（5 步）
1. 初始化项目脚手架
2. 搭建路由/状态管理/API 封装
3. 实现 AES 加解密工具模块
4. 实现登录/注册页面
5. 实现仪表盘及密码管理全套功能


# 工程结构
```code
safebox-web/
├── src/
│   ├── main.ts                    # 入口（注册 Element Plus / Pinia / Router）
│   ├── App.vue                    # 根组件
│   ├── router/index.ts            # 路由 + 导航守卫
│   ├── stores/auth.ts             # 认证状态管理
│   ├── api/
│   │   ├── request.ts             # Axios 实例 + 拦截器
│   │   ├── user.ts                # 用户 API
│   │   └── pwdManage.ts           # 密码管理 API
│   ├── utils/crypto.ts            # AES-256-CBC 加解密工具
│   ├── types/index.ts             # TypeScript 类型定义
│   ├── views/
│   │   ├── Login.vue              # 登录页（渐变背景+玻璃态卡片）
│   │   ├── Register.vue           # 注册页
│   │   └── Dashboard.vue          # 仪表盘主页
│   ├── components/
│   │   ├── PasswordCard.vue       # 密码卡片（脱敏显示 + 解密按钮 + 复制）
│   │   ├── AddPasswordDialog.vue  # 新增密码弹窗
│   │   ├── EncryptDialog.vue      # 加密密钥输入弹窗
│   │   └── DecryptDialog.vue      # 解密密钥输入弹窗
│   └── styles/main.css            # 全局样式（Tailwind + 自定义动画）
├── Dockerfile                     # 多阶段构建：Node 构建 → Nginx 部署
├── nginx.conf                     # Nginx 配置（SPA 路由 + API 反向代理）
├── docker_build.sh                # 构建脚本
├── docker_push.sh                 # 推送镜像脚本
├── vite.config.ts                 # Vite 配置（含开发代理到后端 8888 端口）
├── tailwind.config.js             # Tailwind 配置
└── package.json

```

# 本地启动调试
```bash
# 1. 安装依赖
cd safebox-web && npm install

# 2. 启动开发服务器（自动代理 /api 到 localhost:8888）
npm run dev

# 浏览器访问 http://localhost:3000
```

Vite 已配置代理：前端 `http://localhost:3000/api/*` → 后端 `http://localhost:8888/*`，无需额外配置 CORS

# docker 构建和部署
```bash
# 构建镜像
./docker_build.sh latest

# 或指定仓库推送
./docker_build.sh latest registry.example.com
./docker_push.sh latest registry.example.com

# 运行容器
docker run -d -p 3000:80 safebox-web:latest

# 如果需要与后端一起部署，使用 docker-compose：
# 前端 :3000 (nginx) → 反向代理后端 backend:8888

```
