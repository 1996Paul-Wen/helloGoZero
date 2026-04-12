#!/bin/bash
# 运行容器
docker run -d \
  -p 8888:8888 \
  -v /Users/sauyinman/workspace/git_repo/helloGoZero/safebox/etc/safebox-api.yaml:/root/etc/safebox-api.yaml \
  --name safebox \
  safebox-service