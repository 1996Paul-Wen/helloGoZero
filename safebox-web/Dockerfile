# Stage 1: 构建阶段
FROM node:20-alpine AS builder

WORKDIR /app

# 先复制依赖定义文件，利用 Docker 缓存层加速构建
COPY package.json package-lock.json* ./

RUN npm ci

COPY . .

ARG VITE_API_BASE_URL=/api
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}

RUN npm run build

# Stage 2: 生产运行阶段（Nginx）
FROM nginx:alpine

# 复制构建产物到 Nginx 静态目录
COPY --from=builder /app/dist /usr/share/nginx/html

# 自定义 Nginx 配置：支持 SPA 路由 + API 反向代理
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
