-- User balances and transactions
CREATE TABLE IF NOT EXISTS user_balances (
  user_id VARCHAR(128) NOT NULL,
  point_type_id CHAR(36) NOT NULL,
  balance BIGINT NOT NULL DEFAULT 0,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (user_id, point_type_id),
  INDEX idx_point_type (point_type_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS transactions (
  id CHAR(36) NOT NULL,
  user_id VARCHAR(128) NOT NULL,
  point_type_id CHAR(36) NOT NULL,
  amount BIGINT NOT NULL,
  type ENUM('credit','debit') NOT NULL,
  reason VARCHAR(255) NULL,
  before_balance BIGINT NOT NULL,
  after_balance BIGINT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  INDEX idx_user (user_id, point_type_id, created_at DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


