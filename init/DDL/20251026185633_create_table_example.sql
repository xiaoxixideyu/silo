CREATE TABLE IF NOT EXISTS example
(
    id         BIGINT NOT NULL PRIMARY KEY,
    info       VARCHAR,
    val        BIGINT,
    version    BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

--
-- insert migration record
INSERT INTO sys_migration (name) VALUES ('20251026185633_create_table_example.sql');
