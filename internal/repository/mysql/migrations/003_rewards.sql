-- Reward rules and distributions
CREATE TABLE IF NOT EXISTS reward_rules (
  id CHAR(36) NOT NULL,
  point_type_id CHAR(36) NOT NULL,
  min_rank INT NOT NULL,
  max_rank INT NOT NULL,
  reward_amount BIGINT NOT NULL,
  reward_point_type_id CHAR(36) NOT NULL,
  active TINYINT(1) NOT NULL DEFAULT 1,
  PRIMARY KEY (id),
  INDEX idx_point (point_type_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS reward_distributions (
  id CHAR(36) NOT NULL,
  snapshot_id CHAR(36) NOT NULL,
  executed_at BIGINT NULL,
  status ENUM('pending','completed','failed') NOT NULL DEFAULT 'pending',
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


