# Silo

Silo 是一个用 Go 编写的 HTTP 服务脚手架，旨在简化服务开发流程。

## 项目理念

**Silo** 的含义是 **"Service Is Logic Only"** - 服务只需要关注业务逻辑。

## 特性

使用 Silo 进行服务开发时，以下功能已经集成在脚手架中：

- 🔐 **鉴权** - 完整的身份验证和授权机制
- 📊 **监控** - 内置监控和指标收集
- 📝 **日志** - 结构化日志记录
- 🗄️ **数据库操作** - 数据访问层和 ORM 集成

## 快速开始

使用 Makefile 可以快速初始化和构建项目：

```bash
# 完整的项目构建和初始化
make build
```

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

开发者可以专注于业务逻辑的实现，而无需关心基础设施组件的搭建和配置。

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