#!/bin/bash

# 所有依赖服务都已启动，运行 Go 项目
exec sh -c "cd /go/src && air"

echo "项目成功运行，地址为：http://ip:8080"
