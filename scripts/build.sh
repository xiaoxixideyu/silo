#!/bin/bash

# 获取cmd目录下的所有服务目录
SERVICES=$(find cmd -maxdepth 1 -type d -not -path "cmd" -exec basename {} \;)

# 检查是否指定了服务
if [ -z "$1" ]; then
    echo "Starting complete project build and initialization for all services..."
    echo "Step 1: Running go mod tidy..."
    go mod tidy
    echo "Step 2: Running wire dependency injection for all services..."
    ./scripts/wire.sh
    echo "Step 3: Generating Swagger documentation for all services..."
    ./scripts/swag.sh
    echo "Step 4: Creating build directory..."
    mkdir -p build
    echo "Step 5: Building all services binaries to build directory..."
    for dir in $SERVICES; do
        echo "Building $dir..."
        go build -o build/$dir ./cmd/$dir
    done
    echo "Step 6: Initializing database for all services..."
    for dir in $SERVICES; do
        echo "Initializing database for $dir..."
        ./build/$dir -init -config configs/$dir.yaml
    done
    echo "Project build and initialization completed successfully for all services!"
    echo "Binaries created in build directory: $SERVICES"
else
    SERVICE=$1
    if echo "$SERVICES" | grep -q "\b$SERVICE\b"; then
        echo "Starting project build and initialization for $SERVICE..."
        echo "Step 1: Running go mod tidy..."
        go mod tidy
        echo "Step 2: Running wire dependency injection for $SERVICE..."
        ./scripts/wire.sh $SERVICE
        echo "Step 3: Generating Swagger documentation for $SERVICE..."
        ./scripts/swag.sh $SERVICE
        echo "Step 4: Creating build directory..."
        mkdir -p build
        echo "Step 5: Building $SERVICE binary to build directory..."
        go build -o build/$SERVICE ./cmd/$SERVICE
        echo "Step 6: Initializing database for $SERVICE..."
        ./build/$SERVICE -init -config configs/$SERVICE.yaml
        echo "Project build and initialization completed successfully for $SERVICE!"
        echo "Binary created in build directory: $SERVICE"
    else
        echo "Error: Service '$SERVICE' not found. Available services: $SERVICES"
        exit 1
    fi
fi