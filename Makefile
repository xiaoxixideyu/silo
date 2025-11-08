# Makefile for silo

# 获取cmd目录下的所有服务目录
SERVICES := $(shell find cmd -maxdepth 1 -type d -not -path "cmd" -exec basename {} \;)

.PHONY: help gen-ddl wire swag build

# 默认目标
help:
	@echo "Available targets:"
	@echo "  help    - Show this help message"
	@echo "  gen-ddl - Generate a new DDL migration file"
	@echo "  wire    - Run wire dependency injection for all services or specified service"
	@echo "  swag    - Generate Swagger documentation for all services or specified service"
	@echo "  build   - Complete project initialization and build for all services or specified service"
	@echo ""
	@echo "Available services: $(SERVICES)"
	@echo ""
	@echo "Usage for gen-ddl:"
	@echo "  make gen-ddl filename=<your_filename>"
	@echo "  Example: make gen-ddl filename=create_table_users"
	@echo ""
	@echo "Usage for wire:"
	@echo "  make wire                    # Run wire for all services"
	@echo "  make wire service=admin      # Run wire for specific service"
	@echo ""
	@echo "Usage for swag:"
	@echo "  make swag                    # Generate swagger for all services"
	@echo "  make swag service=admin      # Generate swagger for specific service"
	@echo ""
	@echo "Usage for build:"
	@echo "  make build                   # Build all services"
	@echo "  make build service=admin     # Build specific service"

# 生成 DDL 迁移文件
gen-ddl:
ifndef filename
	@echo "Error: filename parameter is required"
	@echo "Usage: make gen-ddl filename=<your_filename>"
	@exit 1
endif
	@echo "Generating DDL migration file: $(filename)"
	@./scripts/generate_migration.sh init/DDL $(filename)

# 运行 wire 依赖注入
wire:
	@./scripts/wire.sh $(service)

# 生成 Swagger 文档
swag:
	@./scripts/swag.sh $(service)

# 完整的项目构建和初始化
build:
	@./scripts/build.sh $(service)