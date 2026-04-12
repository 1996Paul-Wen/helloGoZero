#!/bin/bash
set -e

IMAGE_NAME="safebox-web"
IMAGE_TAG="${1:-latest}"
REGISTRY="${2:-}"

if [ -z "$REGISTRY" ]; then
  echo "❌ 用法: $0 <tag> <registry>"
  echo "   示例: $0 latest registry.example.com"
  exit 1
fi

FULL_TAG="${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"

echo "📤 推送镜像到仓库..."
echo "   Target: ${FULL_TAG}"

docker push "${FULL_TAG}"

echo "✅ 推送完成！"
