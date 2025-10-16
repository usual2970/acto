-- Point types table
CREATE TABLE IF NOT EXISTS point_types (
  id CHAR(36) NOT NULL,
  name VARCHAR(128) NOT NULL UNIQUE,
  display_name VARCHAR(256) NOT NULL,
  description TEXT NULL,
  enabled TINYINT(1) NOT NULL DEFAULT 1,
  deleted_at BIGINT NULL,
  created_at BIGINT NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


