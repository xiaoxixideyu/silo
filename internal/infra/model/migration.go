package model

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gorm.io/gorm"
)

func Migration(db *gorm.DB) error {
	// 检查 sys_migration 表是否存在
	if err := ensureSysMigrationTable(db); err != nil {
		return fmt.Errorf("failed to ensure sys_migration table: %w", err)
	}

	// 获取所有已执行的迁移
	executedMigrations, err := getExecutedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get executed migrations: %w", err)
	}

	// 获取所有需要执行的迁移文件
	migrationFiles, err := getMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// 执行未执行的迁移
	for _, file := range migrationFiles {
		if !isMigrationExecuted(file.Name(), executedMigrations) {
			if err := executeMigration(db, file); err != nil {
				return fmt.Errorf("failed to execute migration %s: %w", file.Name(), err)
			}
			slog.Info("Migration executed successfully", "file", file.Name())
		}
	}

	return nil
}

// ensureSysMigrationTable 确保 sys_migration 表存在
func ensureSysMigrationTable(db *gorm.DB) error {
	// 检查表是否存在
	var count int64
	err := db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_name = 'sys_migration'").Scan(&count).Error
	if err != nil {
		return err
	}

	// 如果表不存在，则创建
	if count == 0 {
		// 读取并执行创建表的 SQL
		sqlFile := "init/DDL/init_create_table_sys_migration.sql"
		// 尝试从当前工作目录和项目根目录读取文件
		content, err := os.ReadFile(sqlFile)
		if err != nil {
			// 如果从当前目录读取失败，尝试从项目根目录读取
			sqlFile = "../../init/DDL/init_create_table_sys_migration.sql"
			content, err = os.ReadFile(sqlFile)
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %w", sqlFile, err)
			}
		}

		// 先移除注释，再分割 SQL 语句并执行
		cleanContent := removeComments(string(content))
		statements := strings.Split(cleanContent, ";")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt != "" {
				if err := db.Exec(stmt).Error; err != nil {
					return fmt.Errorf("failed to execute SQL: %s, error: %w", stmt, err)
				}
			}
		}

		slog.Info("sys_migration table created successfully")
	}

	return nil
}

// getExecutedMigrations 获取已执行的迁移
func getExecutedMigrations(db *gorm.DB) (map[string]bool, error) {
	var migrations []SysMigration
	err := db.Find(&migrations).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string]bool)
	for _, migration := range migrations {
		result[migration.Name] = true
	}

	return result, nil
}

// getMigrationFiles 获取所有迁移文件
func getMigrationFiles() ([]os.FileInfo, error) {
	dir := "init/DDL"
	files, err := os.ReadDir(dir)
	if err != nil {
		// 如果从当前目录读取失败，尝试从项目根目录读取
		dir = "../../init/DDL"
		files, err = os.ReadDir(dir)
		if err != nil {
			return nil, err
		}
	}

	// 过滤出 SQL 文件并按名称排序
	var sqlFiles []os.FileInfo
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, info)
		}
	}

	// 按文件名排序
	sort.Slice(sqlFiles, func(i, j int) bool {
		return sqlFiles[i].Name() < sqlFiles[j].Name()
	})

	return sqlFiles, nil

}

// isMigrationExecuted 检查迁移是否已执行
func isMigrationExecuted(fileName string, executedMigrations map[string]bool) bool {
	return executedMigrations[fileName]
}

// removeComments 从SQL内容中移除注释
func removeComments(sqlContent string) string {
	lines := strings.Split(sqlContent, "\n")
	var result []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// 跳过以 -- 开头的注释行
		if strings.HasPrefix(trimmed, "--") || trimmed == "" {
			continue
		}
		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

// executeMigration 执行迁移
func executeMigration(db *gorm.DB, file os.FileInfo) error {
	filePath := filepath.Join("init/DDL", file.Name())
	content, err := os.ReadFile(filePath)
	if err != nil {
		// 如果从当前目录读取失败，尝试从项目根目录读取
		filePath = filepath.Join("../../init/DDL", file.Name())
		content, err = os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filePath, err)
		}
	}

	// 先移除注释，再分割 SQL 语句并执行
	cleanContent := removeComments(string(content))
	statements := strings.Split(cleanContent, ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		// 跳过空语句
		if stmt != "" {
			if err := db.Exec(stmt).Error; err != nil {
				// 如果是表已存在的错误，忽略并继续
				if strings.Contains(err.Error(), "already exists") {
					slog.Warn("Table already exists, skipping", "sql", stmt)
					continue
				}
				return fmt.Errorf("failed to execute SQL: %s, error: %w", stmt, err)
			}
		}
	}

	return nil
}
