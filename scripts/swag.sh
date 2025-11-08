#!/bin/bash

# 获取cmd目录下的所有服务目录
SERVICES=$(find cmd -maxdepth 1 -type d -not -path "cmd" -exec basename {} \;)

# 检查是否指定了服务
if [ -z "$1" ]; then
    echo "Generating Swagger documentation for all services..."
    for dir in $SERVICES; do
        echo "Generating Swagger documentation for $dir..."
        cd cmd/$dir && go run github.com/swaggo/swag/cmd/swag@latest init -g ../../internal/interface/$dir/handler.go -o ../../api/$dir --parseDependency --parseInternal -d .,../../internal/interface/$dir,../../internal/infra && cd ../..
    done
    echo "Swagger documentation generated successfully for all services"
else
    SERVICE=$1
    if echo "$SERVICES" | grep -q "\b$SERVICE\b"; then
        echo "Generating Swagger documentation for $SERVICE..."
        cd cmd/$SERVICE && go run github.com/swaggo/swag/cmd/swag@latest init -g ../../internal/interface/$SERVICE/handler.go -o ../../api/$SERVICE --parseDependency --parseInternal -d .,../../internal/interface/$SERVICE,../../internal/infra && cd ../..
        echo "Swagger documentation generated successfully for $SERVICE"
    else
        echo "Error: Service '$SERVICE' not found. Available services: $SERVICES"
        exit 1
    fi
fi