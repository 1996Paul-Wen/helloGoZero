#!/bin/bash
set -e

IMAGE_NAME="safebox-web"
IMAGE_TAG="${1:-latest}"
REGISTRY="${2:-}"

echo "📦 构建 SafeBox Web 前端 Docker 镜像..."
echo "   Image: ${IMAGE_NAME}:${IMAGE_TAG}"

# 获取脚本所在目录的绝对路径
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
echo "使用 $SCRIPT_DIR 路径作为 Docker 构建上下文"

docker build \
  --platform linux/amd64 \
  -t ${IMAGE_NAME}:${IMAGE_TAG} \
  --build-arg VITE_API_BASE_URL=/api \
  "$SCRIPT_DIR"

if [ -n "$REGISTRY" ]; then
  FULL_TAG="${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"
  echo "🏷️  Tagging as ${FULL_TAG}"
  docker tag ${IMAGE_NAME}:${IMAGE_TAG} ${FULL_TAG}
fi

echo "✅ 构建完成！"
echo ""
echo "🚀 运行命令："
# echo "   docker run -d -p 80:80 ${IMAGE_NAME}:${IMAGE_TAG}"

echo "docker run -d \
  -p 80:80 \
  --add-host host.docker.internal:host-gateway \
  -v /Users/sauyinman/workspace/git_repo/helloGoZero/safebox-web/nginx.conf:/etc/nginx/conf.d/default.conf:ro \
  ${IMAGE_NAME}:${IMAGE_TAG}"
