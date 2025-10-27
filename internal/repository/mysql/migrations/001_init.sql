/*
 Navicat Premium Dump SQL

 Source Server         : acto
 Source Server Type    : MySQL
 Source Server Version : 80041 (8.0.41)
 Source Host           : localhost:33061
 Source Schema         : acto

 Target Server Type    : MySQL
 Target Server Version : 80041 (8.0.41)
 File Encoding         : 65001

 Date: 27/10/2025 21:27:05
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for point_types
-- ----------------------------
DROP TABLE IF EXISTS `point_types`;
CREATE TABLE `point_types` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `uri` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `display_name` varchar(256) NOT NULL,
  `description` text,
  `enabled` tinyint(1) NOT NULL DEFAULT '1',
  `deleted_at` bigint DEFAULT NULL,
  `created_at` bigint NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_uri` (`uri`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for redemption_costs
-- ----------------------------
DROP TABLE IF EXISTS `redemption_costs`;
CREATE TABLE `redemption_costs` (
  `reward_id` char(36) NOT NULL,
  `point_type_id` bigint NOT NULL,
  `amount` bigint NOT NULL,
  PRIMARY KEY (`reward_id`,`point_type_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for redemption_records
-- ----------------------------
DROP TABLE IF EXISTS `redemption_records`;
CREATE TABLE `redemption_records` (
  `id` char(36) NOT NULL,
  `user_id` varchar(128) NOT NULL,
  `reward_id` char(36) NOT NULL,
  `created_at` bigint NOT NULL,
  `status` enum('completed','pending','cancelled') NOT NULL DEFAULT 'completed',
  PRIMARY KEY (`id`),
  KEY `idx_user` (`user_id`,`created_at` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for redemption_rewards
-- ----------------------------
DROP TABLE IF EXISTS `redemption_rewards`;
CREATE TABLE `redemption_rewards` (
  `id` char(36) NOT NULL,
  `name` varchar(128) NOT NULL,
  `description` text,
  `quantity` int NOT NULL DEFAULT '0',
  `enabled` tinyint(1) NOT NULL DEFAULT '1',
  `total_redeemed` int NOT NULL DEFAULT '0',
  `created_at` bigint NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for reward_distributions
-- ----------------------------
DROP TABLE IF EXISTS `reward_distributions`;
CREATE TABLE `reward_distributions` (
  `id` char(36) NOT NULL,
  `snapshot_id` char(36) NOT NULL,
  `executed_at` bigint DEFAULT NULL,
  `status` enum('pending','completed','failed') NOT NULL DEFAULT 'pending',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for reward_rules
-- ----------------------------
DROP TABLE IF EXISTS `reward_rules`;
CREATE TABLE `reward_rules` (
  `id` char(36) NOT NULL,
  `point_type_id` bigint NOT NULL,
  `min_rank` int NOT NULL,
  `max_rank` int NOT NULL,
  `reward_amount` bigint NOT NULL,
  `reward_point_type_id` char(36) NOT NULL,
  `active` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`,`point_type_id`) USING BTREE,
  KEY `idx_point` (`point_type_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for transactions
-- ----------------------------
DROP TABLE IF EXISTS `transactions`;
CREATE TABLE `transactions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` varchar(128) NOT NULL,
  `point_type_id` char(36) NOT NULL,
  `amount` bigint NOT NULL,
  `type` enum('credit','debit') NOT NULL,
  `reason` varchar(255) DEFAULT NULL,
  `before_balance` bigint NOT NULL,
  `after_balance` bigint NOT NULL,
  `created_at` bigint NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user` (`user_id`,`point_type_id`,`created_at` DESC)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for user_balances
-- ----------------------------
DROP TABLE IF EXISTS `user_balances`;
CREATE TABLE `user_balances` (
  `user_id` varchar(128) NOT NULL,
  `point_type_id` bigint NOT NULL,
  `balance` bigint NOT NULL DEFAULT '0',
  `updated_at` bigint NOT NULL,
  PRIMARY KEY (`user_id`,`point_type_id`),
  KEY `idx_point_type` (`point_type_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

SET FOREIGN_KEY_CHECKS = 1;