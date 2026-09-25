package dbs

import (
	"database/sql"
	"fmt"
)

// EnsureSchema creates required tables if they do not exist.
func EnsureSchema() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INT NOT NULL AUTO_INCREMENT,
			username VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
			address VARCHAR(255) NOT NULL,
			role VARCHAR(64) NOT NULL DEFAULT 'business',
			PRIMARY KEY (id),
			UNIQUE KEY uk_username (username)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS tx_logs (
			id BIGINT NOT NULL AUTO_INCREMENT,
			action VARCHAR(64) NOT NULL,
			business_id VARCHAR(128) NOT NULL,
			tx_hash VARCHAR(128) NOT NULL,
			status INT NOT NULL,
			message VARCHAR(255) NOT NULL,
			operator_addr VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			KEY idx_business_id (business_id),
			KEY idx_created_at (created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS enterprises (
			enterprise_id VARCHAR(128) NOT NULL,
			name VARCHAR(255) NOT NULL,
			wallet VARCHAR(255) NOT NULL,
			is_core TINYINT(1) NOT NULL DEFAULT 0,
			is_financial TINYINT(1) NOT NULL DEFAULT 0,
			register_time BIGINT NOT NULL DEFAULT 0,
			PRIMARY KEY (enterprise_id),
			UNIQUE KEY uk_enterprise_wallet (wallet)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS receivables (
			receivable_id VARCHAR(128) NOT NULL,
			order_id VARCHAR(128) NOT NULL,
			issuer VARCHAR(255) NOT NULL,
			payer VARCHAR(255) NOT NULL,
			amount DECIMAL(30,0) NOT NULL DEFAULT 0,
			issue_time BIGINT NOT NULL DEFAULT 0,
			due_time BIGINT NOT NULL DEFAULT 0,
			status TINYINT UNSIGNED NOT NULL DEFAULT 0,
			PRIMARY KEY (receivable_id),
			KEY idx_order_id (order_id),
			KEY idx_status (status)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS financing_records (
			id BIGINT NOT NULL AUTO_INCREMENT,
			receivable_id VARCHAR(128) NOT NULL,
			financial_inst VARCHAR(255) NOT NULL,
			amount DECIMAL(30,0) NOT NULL DEFAULT 0,
			interest_rate DECIMAL(30,0) NOT NULL DEFAULT 0,
			financing_time BIGINT NOT NULL DEFAULT 0,
			PRIMARY KEY (id),
			KEY idx_financing_receivable (receivable_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS documents (
			id BIGINT NOT NULL AUTO_INCREMENT,
			name VARCHAR(255) NOT NULL,
			cid VARCHAR(255) NOT NULL,
			receivable_id VARCHAR(128) NOT NULL DEFAULT '',
			order_id VARCHAR(128) NOT NULL DEFAULT '',
			doc_type VARCHAR(64) NOT NULL DEFAULT 'other',
			status VARCHAR(64) NOT NULL DEFAULT 'uploaded',
			source VARCHAR(64) NOT NULL DEFAULT 'ipfs',
			tag VARCHAR(128) NOT NULL DEFAULT '',
			operator VARCHAR(255) NOT NULL DEFAULT '',
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			KEY idx_cid (cid),
			KEY idx_doc_receivable (receivable_id),
			KEY idx_doc_created_at (created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
	}

	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			return fmt.Errorf("schema init failed: %w", err)
		}
	}
	migrations := []struct {
		column string
		def    string
	}{
		{"receivable_id", "ALTER TABLE documents ADD COLUMN receivable_id VARCHAR(128) NOT NULL DEFAULT ''"},
		{"order_id", "ALTER TABLE documents ADD COLUMN order_id VARCHAR(128) NOT NULL DEFAULT ''"},
		{"doc_type", "ALTER TABLE documents ADD COLUMN doc_type VARCHAR(64) NOT NULL DEFAULT 'other'"},
		{"status", "ALTER TABLE documents ADD COLUMN status VARCHAR(64) NOT NULL DEFAULT 'uploaded'"},
		{"source", "ALTER TABLE documents ADD COLUMN source VARCHAR(64) NOT NULL DEFAULT 'ipfs'"},
	}
	for _, item := range migrations {
		exists, err := columnExists("documents", item.column)
		if err != nil {
			return fmt.Errorf("schema inspect failed: %w", err)
		}
		if exists {
			continue
		}
		if _, err = DB.Exec(item.def); err != nil {
			return fmt.Errorf("schema migration failed: %w", err)
		}
	}
	return nil
}

func columnExists(tableName, columnName string) (bool, error) {
	var name string
	err := DB.QueryRow(`
		SELECT COLUMN_NAME
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = ?
		LIMIT 1`, tableName, columnName).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
