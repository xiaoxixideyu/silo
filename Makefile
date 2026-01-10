# Makefile for silo

# 获取cmd目录下的所有服务目录
SERVICES := $(shell find cmd -maxdepth 1 -type d -not -path "cmd" -exec basename {} \;)

.PHONY: help gen-ddl wire swag build deploy-elk cleanup-elk import-dashboard validate-logs docker-up docker-down docker-restart docker-logs docker-ps docker-app-up docker-app-down start-dev-env

# 默认目标
help:
	@echo "Available targets:"
	@echo "  help    - Show this help message"
	@echo "  gen-ddl - Generate a new DDL migration file"
	@echo "  wire    - Run wire dependency injection for all services or specified service"
	@echo "  swag    - Generate Swagger documentation for all services or specified service"
	@echo "  build   - Complete project initialization and build for all services or specified service"
	@echo ""
	@echo "ELK 技术栈命令:"
	@echo "  deploy-elk       - Deploy ELK stack (Elasticsearch, Kibana, Filebeat) to Kubernetes"
	@echo "  cleanup-elk       - Remove ELK stack from Kubernetes"
	@echo "  import-dashboard  - Import pre-configured Kibana dashboard"
	@echo "  validate-logs     - Validate log format compatibility with ELK"
	@echo ""
	@echo "Docker Compose 命令:"
	@echo "  docker-up         - Start all services with Docker Compose"
	@echo "  docker-down       - Stop all services with Docker Compose"
	@echo "  docker-restart    - Restart all services"
	@echo "  docker-logs       - View logs from all services"
	@echo "  docker-ps         - Show status of all services"
	@echo ""
	@echo "Docker Compose 应用服务命令:"
	@echo "  docker-app-up       - Start application services (admin & website)"
	@echo "  docker-app-down     - Stop application services"
	@echo "  start-dev-env       - Start complete development environment (all services)"
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

# 部署 ELK 技术栈
deploy-elk:
	@echo "Deploying ELK stack to Kubernetes..."
	@./scripts/deploy_elk.sh

# 清理 ELK 技术栈
cleanup-elk:
	@echo "Cleaning up ELK stack from Kubernetes..."
	@./scripts/cleanup_elk.sh

# 导入 Kibana 仪表板
import-dashboard:
	@echo "Importing Kibana dashboard..."
	@./scripts/import_dashboard.sh

# 验证日志格式
validate-logs:
	@echo "Validating log format..."
	@./scripts/validate_log_format.sh

# Docker Compose 命令
docker-up:
	@echo "Starting all services with Docker Compose..."
	@docker-compose up -d
	@echo "Waiting for services to be ready..."
	@sleep 5
	@echo ""
	@echo "✅ Services started successfully!"
	@echo ""
	@echo "📋 Service URLs:"
	@echo "  - PostgreSQL: localhost:15432"
	@echo "  - Redis: localhost:6379"
	@echo "  - Elasticsearch: http://localhost:9200"
	@echo "  - Kibana: http://localhost:5601"
	@echo ""
	@echo "💡 Use 'make docker-logs' to view logs"

docker-down:
	@echo "Stopping all services..."
	@docker-compose down
	@echo "✅ Services stopped"

docker-restart:
	@echo "Restarting all services..."
	@docker-compose restart
	@echo "✅ Services restarted"

docker-logs:
	@docker-compose logs -f

docker-ps:
	@docker-compose ps

# Docker Compose 应用服务命令
docker-app-up:
	@echo "Starting application services..."
	@docker-compose -f docker-compose.app.yml up -d
	@echo "✅ Application services started!"
	@echo ""
	@echo "📋 Service URLs:"
	@echo "  - Admin: http://localhost:8080"
	@echo "  - Website: http://localhost:8081"

docker-app-down:
	@echo "Stopping application services..."
	@docker-compose -f docker-compose.app.yml down
	@echo "✅ Application services stopped"

# 一键启动开发环境
start-dev-env:
	@echo "Starting complete development environment..."
	@./scripts/start_dev_env.sh