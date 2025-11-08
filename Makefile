# Makefile for silo

.PHONY: help gen-ddl wire swag build

# 默认目标
help:
	@echo "Available targets:"
	@echo "  help    - Show this help message"
	@echo "  gen-ddl - Generate a new DDL migration file"
	@echo "  wire    - Run wire dependency injection in cmd/admin and cmd/website"
	@echo "  swag    - Generate Swagger documentation for admin and website"
	@echo "  build   - Complete project initialization and build"
	@echo ""
	@echo "Usage for gen-ddl:"
	@echo "  make gen-ddl filename=<your_filename>"
	@echo "  Example: make gen-ddl filename=create_table_users"

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
	@echo "Running wire in cmd/admin..."
	@cd cmd/admin && go run github.com/google/wire/cmd/wire@latest
	@echo "Running wire in cmd/website..."
	@cd cmd/website && go run github.com/google/wire/cmd/wire@latest
	@echo "Wire generation completed for both services"

# 生成 Swagger 文档
swag:
	@echo "Generating Swagger documentation for admin..."
	@cd cmd/admin && go run github.com/swaggo/swag/cmd/swag@latest init -g ../../internal/interface/admin/handler.go -o ../../api/admin --parseDependency --parseInternal -d .,../../internal/interface/admin,../../internal/infra
	@echo "Generating Swagger documentation for website..."
	@cd cmd/website && go run github.com/swaggo/swag/cmd/swag@latest init -g ../../internal/interface/website/handler.go -o ../../api/website --parseDependency --parseInternal -d .,../../internal/interface/website,../../internal/infra
	@echo "Swagger documentation generated successfully for both services"

# 完整的项目构建和初始化
build:
	@echo "Starting complete project build and initialization..."
	@echo "Step 1: Running go mod tidy..."
	@go mod tidy
	@echo "Step 2: Running wire dependency injection..."
	@$(MAKE) wire
	@echo "Step 3: Generating Swagger documentation..."
	@$(MAKE) swag
	@echo "Step 4: Creating build directory..."
	@mkdir -p build
	@echo "Step 5: Building admin and website binaries to build directory..."
	@go build -o build/admin ./cmd/admin
	@go build -o build/website ./cmd/website
	@echo "Step 6: Initializing database for admin..."
	@./build/admin -init -config configs/admin.yaml
	@echo "Step 7: Initializing database for website..."
	@./build/website -init -config configs/website.yaml
	@echo "Project build and initialization completed successfully!"
	@echo "Binaries created in build directory: admin, website"