# Silo

Silo 是一个用 Go 编写的 HTTP 服务脚手架，旨在简化服务开发流程。

## 项目理念

**Silo** 的含义是 **"Silo Is Logic Only"** - Silo只需要关注业务逻辑。

## 特性

使用 Silo 进行服务开发时，以下功能已经集成在脚手架中：

- 🔐 **鉴权** - 完整的身份验证和授权机制
- 📊 **监控** - 内置监控和指标收集
- 📝 **日志** - 结构化日志记录，支持 ELK 技术栈集成
- 🗄️ **数据库操作** - 数据访问层和 ORM 集成
- 🔍 **ELK 集成** - 与 Elasticsearch、Filebeat、Kibana 无缝集成

## 快速开始

使用 Makefile 可以快速初始化和构建项目：

```bash
# 完整的项目构建和初始化
make build
```

### 使用 Docker Compose 启动基础服务

本地开发环境可以使用 Docker Compose 快速启动所有基础服务（PostgreSQL、Redis、Elasticsearch、Kibana）：

```bash
# 一键启动完整开发环境（推荐）
make start-dev-env

# 或分步启动
# 1. 启动基础服务
make docker-up

# 2. 启动应用服务
make docker-app-up

# 查看服务状态
make docker-ps

# 查看服务日志
make docker-logs

# 停止基础服务
make docker-down

# 停止应用服务
make docker-app-down
```

服务访问地址：
- **PostgreSQL**: `localhost:25432` (用户名/密码: silo/silo)
- **Redis**: `localhost:26379`
- **Elasticsearch**: `http://localhost:9200`
- **Kibana**: `http://localhost:5601`
- **Admin 服务**: `http://localhost:8080`
- **Website 服务**: `http://localhost:8081`

数据持久化目录：
- `data/postgres/` - PostgreSQL 数据
- `data/redis/` - Redis 数据
- `data/elasticsearch/` - Elasticsearch 索引数据
- `data/filebeat/` - Filebeat 注册信息
- `logs/` - 应用日志文件

详细文档请参考:
- [Docker Compose 本地开发环境](docs/DOCKER_COMPOSE.md)
- [ELK 日志采集流程说明](docs/ELK_LOG_FLOW.md) - 了解日志如何被采集和展示

这个命令会执行以下步骤：
1. 运行 `go mod tidy` 整理依赖
2. 运行 `wire` 依赖注入
3. 生成 Swagger 文档
4. 构建 admin 和 website 二进制文件
5. 初始化数据库

其他常用命令：

```bash
# 生成新的 DDL 迁移文件
make gen-ddl filename=create_table_users

# 运行 wire 依赖注入
make wire

# 生成 Swagger 文档
make swag

# 查看所有可用命令
make help
```

## ELK 技术栈集成

Silo 框架提供了完整的 ELK 技术栈集成方案，实现日志的集中收集、存储和分析。

### 快速部署

```bash
# 部署 ELK 技术栈（Elasticsearch + Kibana + Filebeat）
make deploy-elk

# 访问 Kibana
kubectl port-forward -n silo svc/kibana 5601:5601

# 导入预配置的仪表板
make import-dashboard
```

### 功能特性

- **自动化日志收集**: Filebeat 作为 Sidecar 容器自动采集应用日志
- **结构化日志**: 基于 JSON 格式的结构化日志，易于查询和分析
- **预配置仪表板**: 包含 API 趋势、错误率、响应时间等关键指标
- **Kubernetes 元数据**: 自动关联 Pod、Node 等集群信息

详细文档请参考: [ELK 集成指南](docs/ELK_INTEGRATION.md)

### 清理 ELK

```bash
# 移除 ELK 技术栈
make cleanup-elk
```

## 项目结构

```
├── cmd/           # 应用程序入口
├── internal/      # 内部包
│   ├── app/       # 应用程序
│   ├── config/    # 配置
│   ├── infra/     # 基础设施层
│   ├── interface/ # 接口层
│   └── service/   # 服务层
├── k8s/           # Kubernetes 配置
├── scripts/       # 脚本工具
└── tools/         # 开发工具