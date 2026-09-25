CREATE TABLE IF NOT EXISTS `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `address` varchar(255) NOT NULL,
  `role` varchar(64) NOT NULL DEFAULT 'business',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `tx_logs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `action` varchar(64) NOT NULL,
  `business_id` varchar(128) NOT NULL,
  `tx_hash` varchar(128) NOT NULL,
  `status` int NOT NULL,
  `message` varchar(255) NOT NULL,
  `operator_addr` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_business_id` (`business_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `enterprises` (
  `enterprise_id` varchar(128) NOT NULL,
  `name` varchar(255) NOT NULL,
  `wallet` varchar(255) NOT NULL,
  `is_core` tinyint(1) NOT NULL DEFAULT 0,
  `is_financial` tinyint(1) NOT NULL DEFAULT 0,
  `register_time` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`enterprise_id`),
  UNIQUE KEY `uk_enterprise_wallet` (`wallet`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `receivables` (
  `receivable_id` varchar(128) NOT NULL,
  `order_id` varchar(128) NOT NULL,
  `issuer` varchar(255) NOT NULL,
  `payer` varchar(255) NOT NULL,
  `amount` decimal(30,0) NOT NULL DEFAULT 0,
  `issue_time` bigint NOT NULL DEFAULT 0,
  `due_time` bigint NOT NULL DEFAULT 0,
  `status` tinyint unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`receivable_id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `financing_records` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `receivable_id` varchar(128) NOT NULL,
  `financial_inst` varchar(255) NOT NULL,
  `amount` decimal(30,0) NOT NULL DEFAULT 0,
  `interest_rate` decimal(30,0) NOT NULL DEFAULT 0,
  `financing_time` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `idx_financing_receivable` (`receivable_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `documents` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `cid` varchar(255) NOT NULL,
  `receivable_id` varchar(128) NOT NULL DEFAULT '',
  `order_id` varchar(128) NOT NULL DEFAULT '',
  `doc_type` varchar(64) NOT NULL DEFAULT 'other',
  `status` varchar(64) NOT NULL DEFAULT 'uploaded',
  `source` varchar(64) NOT NULL DEFAULT 'ipfs',
  `tag` varchar(128) NOT NULL DEFAULT '',
  `operator` varchar(255) NOT NULL DEFAULT '',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_cid` (`cid`),
  KEY `idx_doc_receivable` (`receivable_id`),
  KEY `idx_doc_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
