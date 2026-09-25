package dbs

import (
	"database/sql"
	"fmt"
	"time"
)

// EnterpriseDB represents enterprise record in DB.
type EnterpriseDB struct {
	EnterpriseID string `json:"enterpriseId"`
	Name         string `json:"name"`
	Wallet       string `json:"wallet"`
	IsCore       bool   `json:"isCore"`
	IsFinancial  bool   `json:"isFinancial"`
	RegisterTime int64  `json:"registerTime"`
}

// ReceivableDB represents receivable record in DB.
type ReceivableDB struct {
	ReceivableID string `json:"receivableId"`
	OrderID      string `json:"orderId"`
	Issuer       string `json:"issuer"`
	Payer        string `json:"payer"`
	Amount       string `json:"amount"`
	IssueTime    int64  `json:"issueTime"`
	DueTime      int64  `json:"dueTime"`
	Status       uint8  `json:"status"`
	StatusText   string `json:"statusText"`
}

// FinancingRecordDB represents financing record in DB.
type FinancingRecordDB struct {
	ID            int64  `json:"id"`
	ReceivableID  string `json:"receivableId"`
	FinancialInst string `json:"financialInst"`
	Amount        string `json:"amount"`
	InterestRate  string `json:"interestRate"`
	FinancingTime int64  `json:"financingTime"`
}

// DocumentDB represents file metadata in DB.
type DocumentDB struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	CID          string    `json:"cid"`
	ReceivableID string    `json:"receivableId"`
	OrderID      string    `json:"orderId"`
	DocType      string    `json:"docType"`
	Status       string    `json:"status"`
	Source       string    `json:"source"`
	Tag          string    `json:"tag"`
	Operator     string    `json:"operator"`
	CreatedAt    time.Time `json:"createdAt"`
}

func dbStatusText(status uint8) string {
	switch status {
	case 0:
		return "PENDING"
	case 1:
		return "CONFIRMED"
	case 2:
		return "FINANCED"
	case 3:
		return "SETTLED"
	default:
		return "UNKNOWN"
	}
}

func ListDBEnterprises(limit int) ([]EnterpriseDB, error) {
	if limit <= 0 {
		limit = 500
	}
	rows, err := DB.Query(`SELECT enterprise_id, name, wallet, is_core, is_financial, register_time
	FROM enterprises ORDER BY register_time DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]EnterpriseDB, 0)
	for rows.Next() {
		var x EnterpriseDB
		if err = rows.Scan(&x.EnterpriseID, &x.Name, &x.Wallet, &x.IsCore, &x.IsFinancial, &x.RegisterTime); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

func ListDBReceivables(limit int) ([]ReceivableDB, error) {
	if limit <= 0 {
		limit = 1000
	}
	rows, err := DB.Query(`SELECT receivable_id, order_id, issuer, payer, amount, issue_time, due_time, status
	FROM receivables ORDER BY issue_time DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]ReceivableDB, 0)
	for rows.Next() {
		var x ReceivableDB
		if err = rows.Scan(&x.ReceivableID, &x.OrderID, &x.Issuer, &x.Payer, &x.Amount, &x.IssueTime, &x.DueTime, &x.Status); err != nil {
			return nil, err
		}
		x.StatusText = dbStatusText(x.Status)
		out = append(out, x)
	}
	return out, rows.Err()
}

func GetDBReceivable(receivableID string) (ReceivableDB, error) {
	var x ReceivableDB
	err := DB.QueryRow(`SELECT receivable_id, order_id, issuer, payer, amount, issue_time, due_time, status
	FROM receivables WHERE receivable_id = ?`, receivableID).
		Scan(&x.ReceivableID, &x.OrderID, &x.Issuer, &x.Payer, &x.Amount, &x.IssueTime, &x.DueTime, &x.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return ReceivableDB{}, nil
		}
		return ReceivableDB{}, err
	}
	x.StatusText = dbStatusText(x.Status)
	return x, nil
}

func ListDBFinancingByReceivable(receivableID string) ([]FinancingRecordDB, error) {
	rows, err := DB.Query(`SELECT id, receivable_id, financial_inst, amount, interest_rate, financing_time
	FROM financing_records WHERE receivable_id = ? ORDER BY financing_time DESC`, receivableID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]FinancingRecordDB, 0)
	for rows.Next() {
		var x FinancingRecordDB
		if err = rows.Scan(&x.ID, &x.ReceivableID, &x.FinancialInst, &x.Amount, &x.InterestRate, &x.FinancingTime); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

func ListDBDocuments(limit int) ([]DocumentDB, error) {
	if limit <= 0 {
		limit = 200
	}
	rows, err := DB.Query(`SELECT id, name, cid, receivable_id, order_id, doc_type, status, source, tag, operator, created_at
	FROM documents ORDER BY id DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]DocumentDB, 0)
	for rows.Next() {
		var x DocumentDB
		if err = rows.Scan(&x.ID, &x.Name, &x.CID, &x.ReceivableID, &x.OrderID, &x.DocType, &x.Status, &x.Source, &x.Tag, &x.Operator, &x.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

func UpsertEnterprise(item EnterpriseDB) error {
	_, err := DB.Exec(`INSERT INTO enterprises (enterprise_id, name, wallet, is_core, is_financial, register_time)
	VALUES (?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE
		name = VALUES(name),
		wallet = VALUES(wallet),
		is_core = VALUES(is_core),
		is_financial = VALUES(is_financial),
		register_time = VALUES(register_time)`,
		item.EnterpriseID, item.Name, item.Wallet, item.IsCore, item.IsFinancial, item.RegisterTime)
	return err
}

func UpsertReceivable(item ReceivableDB) error {
	_, err := DB.Exec(`INSERT INTO receivables (receivable_id, order_id, issuer, payer, amount, issue_time, due_time, status)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE
		order_id = VALUES(order_id),
		issuer = VALUES(issuer),
		payer = VALUES(payer),
		amount = VALUES(amount),
		issue_time = VALUES(issue_time),
		due_time = VALUES(due_time),
		status = VALUES(status)`,
		item.ReceivableID, item.OrderID, item.Issuer, item.Payer, item.Amount, item.IssueTime, item.DueTime, item.Status)
	return err
}

func InsertFinancingRecord(item FinancingRecordDB) error {
	_, err := DB.Exec(`INSERT INTO financing_records (receivable_id, financial_inst, amount, interest_rate, financing_time)
	VALUES (?, ?, ?, ?, ?)`,
		item.ReceivableID, item.FinancialInst, item.Amount, item.InterestRate, item.FinancingTime)
	return err
}

func InsertDocument(item DocumentDB) error {
	_, err := DB.Exec(`INSERT INTO documents (name, cid, receivable_id, order_id, doc_type, status, source, tag, operator)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		item.Name, item.CID, item.ReceivableID, item.OrderID, item.DocType, item.Status, item.Source, item.Tag, item.Operator)
	return err
}

func GetDBOverview() (map[string]any, error) {
	var userCount int
	if err := DB.QueryRow(`SELECT COUNT(1) FROM users`).Scan(&userCount); err != nil {
		return nil, err
	}

	var receivableCount int
	if err := DB.QueryRow(`SELECT COUNT(1) FROM receivables`).Scan(&receivableCount); err != nil {
		return nil, err
	}

	statusCounters := map[string]int{
		"PENDING":   0,
		"CONFIRMED": 0,
		"FINANCED":  0,
		"SETTLED":   0,
	}
	rows, err := DB.Query(`SELECT status, COUNT(1) AS c FROM receivables GROUP BY status`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var s uint8
		var c int
		if err = rows.Scan(&s, &c); err != nil {
			_ = rows.Close()
			return nil, err
		}
		statusCounters[dbStatusText(s)] = c
	}
	_ = rows.Close()

	var totalAmount sql.NullString
	if err = DB.QueryRow(`SELECT CAST(IFNULL(SUM(amount), 0) AS CHAR) FROM receivables`).Scan(&totalAmount); err != nil {
		return nil, err
	}

	txLogs, err := ListRecentTxLogs(30)
	if err != nil {
		return nil, fmt.Errorf("load tx logs failed: %w", err)
	}

	return map[string]any{
		"users":          userCount,
		"receivableSize": receivableCount,
		"statusCounters": statusCounters,
		"totalAmount":    totalAmount.String,
		"recentTxLogs":   txLogs,
	}, nil
}
