#!/bin/bash

# 1.tag image
# docker tag {imageid} crpi-d9ey6ba0ii5v9rr5.cn-hangzhou.personal.cr.aliyuncs.com/paulwen/safebox-service:latest

# 2.login to aliyun private registry
# docker login --username=nick3335561944 crpi-d9ey6ba0ii5v9rr5.cn-hangzhou.personal.cr.aliyuncs.com

# 3.push image
docker push crpi-d9ey6ba0ii5v9rr5.cn-hangzhou.personal.cr.aliyuncs.com/paulwen/safebox-service:latest