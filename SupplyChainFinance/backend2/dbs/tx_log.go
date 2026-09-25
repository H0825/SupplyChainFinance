package dbs

import (
	"database/sql"
	"time"
)

// TxLog stores chain transaction metadata.
type TxLog struct {
	ID           int64     `json:"id"`
	Action       string    `json:"action"`
	BusinessID   string    `json:"businessId"`
	TxHash       string    `json:"txHash"`
	Status       int64     `json:"status"`
	Message      string    `json:"message"`
	OperatorAddr string    `json:"operatorAddr"`
	CreatedAt    time.Time `json:"createdAt"`
}

func SaveTxLog(log TxLog) error {
	q := `INSERT INTO tx_logs (action, business_id, tx_hash, status, message, operator_addr)
	VALUES (?, ?, ?, ?, ?, ?)`
	_, err := DB.Exec(q, log.Action, log.BusinessID, log.TxHash, log.Status, log.Message, log.OperatorAddr)
	return err
}

func ListRecentTxLogs(limit int) ([]TxLog, error) {
	if limit <= 0 {
		limit = 10
	}
	rows, err := DB.Query(`SELECT id, action, business_id, tx_hash, status, message, operator_addr, created_at
	FROM tx_logs ORDER BY id DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := make([]TxLog, 0, limit)
	for rows.Next() {
		var x TxLog
		if err = rows.Scan(&x.ID, &x.Action, &x.BusinessID, &x.TxHash, &x.Status, &x.Message, &x.OperatorAddr, &x.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, x)
	}
	return logs, rows.Err()
}

func CountTxLogs() (int, error) {
	var total int
	err := DB.QueryRow(`SELECT COUNT(1) FROM tx_logs`).Scan(&total)
	return total, err
}

func ListTxLogsPage(page, pageSize int) ([]TxLog, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 30
	}
	offset := (page - 1) * pageSize
	rows, err := DB.Query(`SELECT id, action, business_id, tx_hash, status, message, operator_addr, created_at
	FROM tx_logs ORDER BY id DESC LIMIT ? OFFSET ?`, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := make([]TxLog, 0, pageSize)
	for rows.Next() {
		var x TxLog
		if err = rows.Scan(&x.ID, &x.Action, &x.BusinessID, &x.TxHash, &x.Status, &x.Message, &x.OperatorAddr, &x.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, x)
	}
	return logs, rows.Err()
}

func GetTxLogByHash(txHash string) (TxLog, error) {
	var x TxLog
	err := DB.QueryRow(`SELECT id, action, business_id, tx_hash, status, message, operator_addr, created_at
	FROM tx_logs WHERE tx_hash = ? ORDER BY id DESC LIMIT 1`, txHash).
		Scan(&x.ID, &x.Action, &x.BusinessID, &x.TxHash, &x.Status, &x.Message, &x.OperatorAddr, &x.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return TxLog{}, nil
		}
		return TxLog{}, err
	}
	return x, nil
}
