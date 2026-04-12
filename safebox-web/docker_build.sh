#!/bin/bash
set -e

IMAGE_NAME="safebox-web"
IMAGE_TAG="${1:-latest}"
REGISTRY="${2:-}"

echo "📦 构建 SafeBox Web 前端 Docker 镜像..."
echo "   Image: ${IMAGE_NAME}:${IMAGE_TAG}"

docker build \
  -t ${IMAGE_NAME}:${IMAGE_TAG} \
  --build-arg VITE_API_BASE_URL=/api \
  "$(cd "$(dirname "$0")" && pwd)"

if [ -n "$REGISTRY" ]; then
  FULL_TAG="${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"
  echo "🏷️  Tagging as ${FULL_TAG}"
  docker tag ${IMAGE_NAME}:${IMAGE_TAG} ${FULL_TAG}
fi

echo "✅ 构建完成！"
echo ""
echo "🚀 运行命令："
echo "   docker run -d -p 3000:80 ${IMAGE_NAME}:${IMAGE_TAG}"
