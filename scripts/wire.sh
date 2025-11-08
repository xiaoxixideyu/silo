#!/bin/bash

# 获取cmd目录下的所有服务目录
SERVICES=$(find cmd -maxdepth 1 -type d -not -path "cmd" -exec basename {} \;)

# 检查是否指定了服务
if [ -z "$1" ]; then
    echo "Running wire for all services..."
    for dir in $SERVICES; do
        echo "Running wire in cmd/$dir..."
        cd cmd/$dir && go run github.com/google/wire/cmd/wire@latest && cd ../..
    done
    echo "Wire generation completed for all services"
else
    SERVICE=$1
    if echo "$SERVICES" | grep -q "\b$SERVICE\b"; then
        echo "Running wire in cmd/$SERVICE..."
        cd cmd/$SERVICE && go run github.com/google/wire/cmd/wire@latest && cd ../..
        echo "Wire generation completed for $SERVICE"
    else
        echo "Error: Service '$SERVICE' not found. Available services: $SERVICES"
        exit 1
    fi
fi