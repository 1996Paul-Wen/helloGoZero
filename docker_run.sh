#!/bin/bash
# 镜像构建之后，运行容器
docker run -d \
  -p 8888:8888 \
  --add-host host.docker.internal:host-gateway \
  -v /Users/sauyinman/workspace/git_repo/helloGoZero/safebox/etc/safebox-api.yaml:/root/etc/safebox-api.yaml \
  --name safebox \
  safebox-service