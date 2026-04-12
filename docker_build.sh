#!/bin/bash
docker build --platform linux/arm64 -t safebox-service .  # 构建 arm64 镜像

# # 构建 amd64 镜像
# docker build --platform linux/amd64 -t crpi-d9ey6ba0ii5v9rr5.cn-hangzhou.personal.cr.aliyuncs.com/paulwen/safebox-service:latest .