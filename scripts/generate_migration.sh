#!/bin/bash

# 检查参数数量
if [ $# -ne 2 ]; then
    echo "Usage: $0 <directory> <filename>"
    echo "Example: $0 init/DDL create_table_users"
    exit 1
fi

# 获取参数
TARGET_DIR="$1"
FILENAME="$2"

# 检查目录是否存在，如果不存在则创建
if [ ! -d "$TARGET_DIR" ]; then
    echo "Directory $TARGET_DIR does not exist. Creating it..."
    mkdir -p "$TARGET_DIR"
fi

# 生成时间戳（格式：YYYYMMDDHHMMSS）
TIMESTAMP=$(date +"%Y%m%d%H%M%S")

# 构建完整的文件名
FULL_FILENAME="${TIMESTAMP}_${FILENAME}.sql"

# 构建完整的文件路径
FILE_PATH="${TARGET_DIR}/${FULL_FILENAME}"

# 生成 SQL 文件内容
cat > "$FILE_PATH" << EOF
-- TODO add your scripts here

--
-- insert migration record
INSERT INTO sys_migration (name) VALUES ('${FULL_FILENAME}');
EOF

echo "Generated migration file: $FILE_PATH"