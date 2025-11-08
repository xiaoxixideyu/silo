CREATE TABLE IF NOT EXISTS sys_migration (
     name VARCHAR NOT NULL PRIMARY KEY,
     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--
-- insert migration record
INSERT INTO sys_migration (name) VALUES ('init_create_table_sys_migration.sql');