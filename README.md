[中文版](README_cn.md) | English

# Silo

Silo is an HTTP service scaffold written in Go, designed to simplify the service development process.

## Project Philosophy

The meaning of **Silo** is **"Silo Is Logic Only"** - Silo only need to focus on business logic.

## Features

When using Silo for service development, the following features are integrated into the scaffold:

- 🔐 **Authentication** - Complete authentication and authorization mechanisms
- 📊 **Monitoring** - Built-in monitoring and metrics collection
- 📝 **Logging** - Structured logging
- 🗄️ **Database Operations** - Data access layer and ORM integration

## Quick Start

Use Makefile to quickly initialize and build the project:

```bash
# Complete project build and initialization
make build
```

This command executes the following steps:
1. Run `go mod tidy` to organize dependencies
2. Run `wire` for dependency injection
3. Generate Swagger documentation
4. Build admin and website binaries
5. Initialize database

Other common commands:

```bash
# Generate new DDL migration file
make gen-ddl filename=create_table_users

# Run wire dependency injection
make wire

# Generate Swagger documentation
make swag

# View all available commands
make help
```

Developers can focus on implementing business logic without worrying about infrastructure component setup and configuration.

## Project Structure

```
├── cmd/           # Application entry points
├── internal/      # Internal packages
│   ├── app/       # Application layer
│   ├── config/    # Configuration
│   ├── infra/     # Infrastructure layer
│   ├── interface/ # Interface layer
│   └── service/   # Service layer
├── k8s/           # Kubernetes configurations
├── scripts/       # Script tools
└── tools/         # Development tools