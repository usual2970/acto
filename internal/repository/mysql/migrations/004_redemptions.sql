-- Redemption rewards and records
CREATE TABLE IF NOT EXISTS redemption_rewards (
  id CHAR(36) NOT NULL,
  name VARCHAR(128) NOT NULL,
  description TEXT NULL,
  quantity INT NOT NULL DEFAULT 0,
  enabled TINYINT(1) NOT NULL DEFAULT 1,
  total_redeemed INT NOT NULL DEFAULT 0,
  created_at BIGINT NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Simple costs table: multiple point types per reward
CREATE TABLE IF NOT EXISTS redemption_costs (
  reward_id CHAR(36) NOT NULL,
  point_type_id CHAR(36) NOT NULL,
  amount BIGINT NOT NULL,
  PRIMARY KEY (reward_id, point_type_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS redemption_records (
  id CHAR(36) NOT NULL,
  user_id VARCHAR(128) NOT NULL,
  reward_id CHAR(36) NOT NULL,
  created_at BIGINT NOT NULL,
  status ENUM('completed','pending','cancelled') NOT NULL DEFAULT 'completed',
  PRIMARY KEY (id),
  INDEX idx_user (user_id, created_at DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


