package router

import (
	"backend/SupplyChainFinance"
	"backend/appconfig"
	"backend/dbs"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/FISCO-BCOS/go-sdk/client"
	"github.com/FISCO-BCOS/go-sdk/conf"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

var (
	Configs         []conf.Config
	SupplyClient    *client.Client
	ContractAddress = common.Address{}
	SupplyInstance  *SupplyChainFinance.SupplyChainFinance
)

func init() {
	var err error
	appCfg, err := appconfig.Load("config.toml")
	if err != nil {
		log.Fatalf("Load app config failed, err: %v", err)
	}
	if !common.IsHexAddress(appCfg.Contract.Address) {
		log.Fatalf("invalid contract address in config.toml: %s", appCfg.Contract.Address)
	}
	ContractAddress = common.HexToAddress(appCfg.Contract.Address)

	Configs, err = conf.ParseConfigFile("config.toml")
	if err != nil {
		log.Fatalf("ParseConfigFile failed, err: %v", err)
	}
	SupplyClient, err = client.Dial(&Configs[0])
	if err != nil {
		log.Fatalf("client.Dial failed, err: %v", err)
	}
	SupplyInstance, err = SupplyChainFinance.NewSupplyChainFinance(ContractAddress, SupplyClient)
	if err != nil {
		log.Fatalf("NewInstance failed, err: %v", err)
	}
}

type txResult struct {
	Hash    string `json:"hash"`
	Status  int64  `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type enterpriseView struct {
	EnterpriseID string `json:"enterpriseId"`
	Name         string `json:"name"`
	Wallet       string `json:"wallet"`
	IsCore       bool   `json:"isCore"`
	IsFinancial  bool   `json:"isFinancial"`
	RegisterTime string `json:"registerTime"`
	Username     string `json:"username,omitempty"`
	Role         string `json:"role,omitempty"`
}

type financingView struct {
	FinancialInst string `json:"financialInst"`
	Amount        string `json:"amount"`
	InterestRate  string `json:"interestRate"`
	FinancingTime string `json:"financingTime"`
}

type receivableView struct {
	ReceivableID string `json:"receivableId"`
	OrderID      string `json:"orderId"`
	Issuer       string `json:"issuer"`
	Payer        string `json:"payer"`
	Amount       string `json:"amount"`
	IssueTime    string `json:"issueTime"`
	DueTime      string `json:"dueTime"`
	Status       uint8  `json:"status"`
	StatusText   string `json:"statusText"`
}

func statusText(status uint8) string {
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

func getOperatorAddr() string {
	opts := SupplyClient.GetTransactOpts()
	if opts == nil {
		return ""
	}
	return opts.From.Hex()
}

func saveTxLog(action, businessID string, result txResult) {
	_ = dbs.SaveTxLog(dbs.TxLog{
		Action:       action,
		BusinessID:   businessID,
		TxHash:       result.Hash,
		Status:       result.Status,
		Message:      result.Message,
		OperatorAddr: getOperatorAddr(),
	})
}

func okTx(hash string, status int64, message string) txResult {
	return txResult{Hash: hash, Status: status, Message: message}
}

func failTx(hash string, status int64, message string, err error) txResult {
	r := txResult{Hash: hash, Status: status, Message: message}
	if err != nil {
		r.Error = err.Error()
	}
	return r
}

func enterpriseToView(e SupplyChainFinance.SupplyChainFinanceEnterprise) enterpriseView {
	registerTime := "0"
	if e.RegisterTime != nil {
		registerTime = e.RegisterTime.String()
	}
	return enterpriseView{
		EnterpriseID: e.EnterpriseId,
		Name:         e.Name,
		Wallet:       e.Wallet.Hex(),
		IsCore:       e.IsCore,
		IsFinancial:  e.IsFinancial,
		RegisterTime: registerTime,
	}
}

func receivableToView(r SupplyChainFinance.SupplyChainFinanceReceivable) receivableView {
	return receivableView{
		ReceivableID: r.ReceivableId,
		OrderID:      r.OrderId,
		Issuer:       r.Issuer.Hex(),
		Payer:        r.Payer.Hex(),
		Amount:       r.Amount.String(),
		IssueTime:    r.IssueTime.String(),
		DueTime:      r.DueTime.String(),
		Status:       r.Status,
		StatusText:   statusText(r.Status),
	}
}

func financingToView(x SupplyChainFinance.SupplyChainFinanceFinancingRecord) financingView {
	return financingView{
		FinancialInst: x.FinancialInst.Hex(),
		Amount:        x.Amount.String(),
		InterestRate:  x.InterestRate.String(),
		FinancingTime: x.FinancingTime.String(),
	}
}

func syncEnterpriseToDBByWallet(wallet common.Address) error {
	e, err := SupplyInstance.GetEnterprise(nil, wallet)
	if err != nil || e.EnterpriseId == "" {
		if err != nil {
			return err
		}
		return fmt.Errorf("enterprise not found on chain")
	}
	registerTime := int64(0)
	if e.RegisterTime != nil {
		registerTime = e.RegisterTime.Int64()
	}
	return dbs.UpsertEnterprise(dbs.EnterpriseDB{
		EnterpriseID: e.EnterpriseId,
		Name:         e.Name,
		Wallet:       e.Wallet.Hex(),
		IsCore:       e.IsCore,
		IsFinancial:  e.IsFinancial,
		RegisterTime: registerTime,
	})
}

func syncReceivableToDB(receivableID string) error {
	r, err := SupplyInstance.GetReceivable(nil, receivableID)
	if err != nil || r.ReceivableId == "" {
		if err != nil {
			return err
		}
		return fmt.Errorf("receivable not found on chain")
	}
	issueTime := int64(0)
	if r.IssueTime != nil {
		issueTime = r.IssueTime.Int64()
	}
	dueTime := int64(0)
	if r.DueTime != nil {
		dueTime = r.DueTime.Int64()
	}
	return dbs.UpsertReceivable(dbs.ReceivableDB{
		ReceivableID: r.ReceivableId,
		OrderID:      r.OrderId,
		Issuer:       r.Issuer.Hex(),
		Payer:        r.Payer.Hex(),
		Amount:       r.Amount.String(),
		IssueTime:    issueTime,
		DueTime:      dueTime,
		Status:       r.Status,
		StatusText:   statusText(r.Status),
	})
}

func syncLatestFinancingRecordToDB(receivableID string) error {
	recs, err := SupplyInstance.GetFinancingRecords(nil, receivableID)
	if err != nil || len(recs) == 0 {
		if err != nil {
			return err
		}
		return fmt.Errorf("financing record not found on chain")
	}
	last := recs[len(recs)-1]
	financingTime := int64(0)
	if last.FinancingTime != nil {
		financingTime = last.FinancingTime.Int64()
	}
	return dbs.InsertFinancingRecord(dbs.FinancingRecordDB{
		ReceivableID:  receivableID,
		FinancialInst: last.FinancialInst.Hex(),
		Amount:        last.Amount.String(),
		InterestRate:  last.InterestRate.String(),
		FinancingTime: financingTime,
	})
}

// RegisterEnterprise registers enterprise on chain.
func RegisterEnterprise(c *gin.Context) {
	var req struct {
		EnterpriseID string `json:"enterpriseId" binding:"required"`
		Name         string `json:"name" binding:"required"`
		IsCore       bool   `json:"isCore"`
		IsFinance    bool   `json:"isFinance"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	operatorAddr := getOperatorAddr()
	if common.IsHexAddress(operatorAddr) {
		existing, err := SupplyInstance.GetEnterprise(nil, common.HexToAddress(operatorAddr))
		if err == nil && existing.EnterpriseId != "" {
			registerTime := int64(0)
			if existing.RegisterTime != nil {
				registerTime = existing.RegisterTime.Int64()
			}
			if syncErr := dbs.UpsertEnterprise(dbs.EnterpriseDB{
				EnterpriseID: existing.EnterpriseId,
				Name:         existing.Name,
				Wallet:       existing.Wallet.Hex(),
				IsCore:       existing.IsCore,
				IsFinancial:  existing.IsFinancial,
				RegisterTime: registerTime,
			}); syncErr != nil {
				log.Printf("sync existing enterprise to db failed: %v", syncErr)
			}
			c.JSON(http.StatusOK, gin.H{
				"message":           "enterprise already registered",
				"alreadyRegistered": true,
				"data":              enterpriseToView(existing),
			})
			return
		}
	}

	tx, receipt, err := SupplyInstance.RegisterEnterprise(
		SupplyClient.GetTransactOpts(),
		req.EnterpriseID,
		req.Name,
		req.IsCore,
		req.IsFinance,
	)
	if err != nil {
		res := failTx("", -1, "register enterprise failed", err)
		saveTxLog("register_enterprise", req.EnterpriseID, res)
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	if receipt.Status == 22 {
		res := failTx(tx.Hash().Hex(), int64(receipt.Status), "register enterprise reverted", fmt.Errorf(receipt.GetErrorMessage()))
		saveTxLog("register_enterprise", req.EnterpriseID, res)
		c.JSON(http.StatusConflict, res)
		return
	}

	res := okTx(tx.Hash().Hex(), int64(receipt.Status), "register enterprise success")
	saveTxLog("register_enterprise", req.EnterpriseID, res)
	if common.IsHexAddress(operatorAddr) {
		if syncErr := dbs.UpsertEnterprise(dbs.EnterpriseDB{
			EnterpriseID: req.EnterpriseID,
			Name:         req.Name,
			Wallet:       operatorAddr,
			IsCore:       req.IsCore,
			IsFinancial:  req.IsFinance,
			RegisterTime: time.Now().Unix(),
		}); syncErr != nil {
			log.Printf("direct enterprise upsert failed: %v", syncErr)
			if fallbackErr := syncEnterpriseToDBByWallet(common.HexToAddress(operatorAddr)); fallbackErr != nil {
				log.Printf("chain enterprise sync fallback failed: %v", fallbackErr)
			}
		}
	}
	c.JSON(http.StatusOK, res)
}

// GetEnterpriseByAddress gets enterprise info by address.
func GetEnterpriseByAddress(c *gin.Context) {
	addrHex := strings.TrimSpace(c.Param("address"))
	if !common.IsHexAddress(addrHex) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid address"})
		return
	}

	e, err := SupplyInstance.GetEnterprise(nil, common.HexToAddress(addrHex))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if e.EnterpriseId == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "enterprise not found on-chain"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": enterpriseToView(e)})
}

// ListEnterprises lists enterprises by known user addresses.
func ListEnterprises(c *gin.Context) {
	users, err := dbs.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	list := make([]enterpriseView, 0)
	for _, u := range users {
		if !common.IsHexAddress(u.Address) {
			continue
		}
		e, callErr := SupplyInstance.GetEnterprise(nil, common.HexToAddress(u.Address))
		if callErr != nil || e.EnterpriseId == "" {
			continue
		}
		item := enterpriseToView(e)
		item.Username = u.Username
		item.Role = u.Role
		list = append(list, item)
	}
	c.JSON(http.StatusOK, gin.H{"count": len(list), "data": list})
}

// IssueReceivable creates receivable on chain.
func IssueReceivable(c *gin.Context) {
	var req struct {
		ReceivableID string `json:"receivableId" binding:"required"`
		OrderID      string `json:"orderId" binding:"required"`
		Payer        string `json:"payer" binding:"required"`
		Amount       int64  `json:"amount" binding:"required"`
		DueTimestamp int64  `json:"dueTimestamp" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !common.IsHexAddress(req.Payer) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payer address"})
		return
	}
	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be > 0"})
		return
	}

	tx, receipt, err := SupplyInstance.IssueReceivable(
		SupplyClient.GetTransactOpts(),
		req.ReceivableID,
		req.OrderID,
		common.HexToAddress(req.Payer),
		big.NewInt(req.Amount),
		big.NewInt(req.DueTimestamp),
	)
	if err != nil {
		res := failTx("", -1, "issue receivable failed", err)
		saveTxLog("issue_receivable", req.ReceivableID, res)
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	if receipt.Status == 22 {
		res := failTx(tx.Hash().Hex(), int64(receipt.Status), "issue receivable reverted", fmt.Errorf(receipt.GetErrorMessage()))
		saveTxLog("issue_receivable", req.ReceivableID, res)
		c.JSON(http.StatusConflict, res)
		return
	}
	res := okTx(tx.Hash().Hex(), int64(receipt.Status), "issue receivable success")
	saveTxLog("issue_receivable", req.ReceivableID, res)
	issuer := getOperatorAddr()
	if syncErr := dbs.UpsertReceivable(dbs.ReceivableDB{
		ReceivableID: req.ReceivableID,
		OrderID:      req.OrderID,
		Issuer:       issuer,
		Payer:        req.Payer,
		Amount:       strconv.FormatInt(req.Amount, 10),
		IssueTime:    time.Now().Unix(),
		DueTime:      req.DueTimestamp,
		Status:       0,
		StatusText:   "PENDING",
	}); syncErr != nil {
		log.Printf("direct receivable upsert failed: %v", syncErr)
		if fallbackErr := syncReceivableToDB(req.ReceivableID); fallbackErr != nil {
			log.Printf("chain receivable sync fallback failed: %v", fallbackErr)
		}
	}
	c.JSON(http.StatusOK, res)
}

// ConfirmReceivable confirms receivable.
func ConfirmReceivable(c *gin.Context) {
	var req struct {
		ReceivableID string `json:"receivableId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tx, receipt, err := SupplyInstance.ConfirmReceivable(SupplyClient.GetTransactOpts(), req.ReceivableID)
	if err != nil {
		res := failTx("", -1, "confirm receivable failed", err)
		saveTxLog("confirm_receivable", req.ReceivableID, res)
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	if receipt.Status == 22 {
		res := failTx(tx.Hash().Hex(), int64(receipt.Status), "confirm receivable reverted", fmt.Errorf(receipt.GetErrorMessage()))
		saveTxLog("confirm_receivable", req.ReceivableID, res)
		c.JSON(http.StatusConflict, res)
		return
	}
	res := okTx(tx.Hash().Hex(), int64(receipt.Status), "confirm receivable success")
	saveTxLog("confirm_receivable", req.ReceivableID, res)
	rec, getErr := dbs.GetDBReceivable(req.ReceivableID)
	if getErr == nil && rec.ReceivableID != "" {
		rec.Status = 1
		rec.StatusText = "CONFIRMED"
		if syncErr := dbs.UpsertReceivable(rec); syncErr != nil {
			log.Printf("direct confirm receivable upsert failed: %v", syncErr)
			if fallbackErr := syncReceivableToDB(req.ReceivableID); fallbackErr != nil {
				log.Printf("chain confirm sync fallback failed: %v", fallbackErr)
			}
		}
	} else if fallbackErr := syncReceivableToDB(req.ReceivableID); fallbackErr != nil {
		log.Printf("chain confirm sync fallback failed: %v", fallbackErr)
	}
	c.JSON(http.StatusOK, res)
}

// RecordFinancing records financing.
func RecordFinancing(c *gin.Context) {
	var req struct {
		ReceivableID string `json:"receivableId" binding:"required"`
		Amount       int64  `json:"amount" binding:"required"`
		InterestRate int64  `json:"interestRate" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be > 0"})
		return
	}
	if req.InterestRate < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "interestRate must be >= 0"})
		return
	}

	tx, receipt, err := SupplyInstance.RecordFinancing(
		SupplyClient.GetTransactOpts(),
		req.ReceivableID,
		big.NewInt(req.Amount),
		big.NewInt(req.InterestRate),
	)
	if err != nil {
		res := failTx("", -1, "record financing failed", err)
		saveTxLog("record_financing", req.ReceivableID, res)
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	if receipt.Status == 22 {
		res := failTx(tx.Hash().Hex(), int64(receipt.Status), "record financing reverted", fmt.Errorf(receipt.GetErrorMessage()))
		saveTxLog("record_financing", req.ReceivableID, res)
		c.JSON(http.StatusConflict, res)
		return
	}
	res := okTx(tx.Hash().Hex(), int64(receipt.Status), "record financing success")
	saveTxLog("record_financing", req.ReceivableID, res)
	rec, getErr := dbs.GetDBReceivable(req.ReceivableID)
	if getErr == nil && rec.ReceivableID != "" {
		rec.Status = 2
		rec.StatusText = "FINANCED"
		if syncErr := dbs.UpsertReceivable(rec); syncErr != nil {
			log.Printf("direct financing receivable upsert failed: %v", syncErr)
		}
	}
	if syncErr := dbs.InsertFinancingRecord(dbs.FinancingRecordDB{
		ReceivableID:  req.ReceivableID,
		FinancialInst: getOperatorAddr(),
		Amount:        strconv.FormatInt(req.Amount, 10),
		InterestRate:  strconv.FormatInt(req.InterestRate, 10),
		FinancingTime: time.Now().Unix(),
	}); syncErr != nil {
		log.Printf("direct financing record insert failed: %v", syncErr)
		if fallbackErr := syncLatestFinancingRecordToDB(req.ReceivableID); fallbackErr != nil {
			log.Printf("chain financing sync fallback failed: %v", fallbackErr)
		}
	}
	if fallbackErr := syncReceivableToDB(req.ReceivableID); fallbackErr != nil {
		log.Printf("chain financing receivable sync fallback failed: %v", fallbackErr)
	}
	c.JSON(http.StatusOK, res)
}

// RecordSettlement marks receivable settled.
func RecordSettlement(c *gin.Context) {
	var req struct {
		ReceivableID string `json:"receivableId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tx, receipt, err := SupplyInstance.RecordSettlement(SupplyClient.GetTransactOpts(), req.ReceivableID)
	if err != nil {
		res := failTx("", -1, "record settlement failed", err)
		saveTxLog("record_settlement", req.ReceivableID, res)
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	if receipt.Status == 22 {
		res := failTx(tx.Hash().Hex(), int64(receipt.Status), "record settlement reverted", fmt.Errorf(receipt.GetErrorMessage()))
		saveTxLog("record_settlement", req.ReceivableID, res)
		c.JSON(http.StatusConflict, res)
		return
	}
	res := okTx(tx.Hash().Hex(), int64(receipt.Status), "record settlement success")
	saveTxLog("record_settlement", req.ReceivableID, res)
	rec, getErr := dbs.GetDBReceivable(req.ReceivableID)
	if getErr == nil && rec.ReceivableID != "" {
		rec.Status = 3
		rec.StatusText = "SETTLED"
		if syncErr := dbs.UpsertReceivable(rec); syncErr != nil {
			log.Printf("direct settlement receivable upsert failed: %v", syncErr)
			if fallbackErr := syncReceivableToDB(req.ReceivableID); fallbackErr != nil {
				log.Printf("chain settlement sync fallback failed: %v", fallbackErr)
			}
		}
	} else if fallbackErr := syncReceivableToDB(req.ReceivableID); fallbackErr != nil {
		log.Printf("chain settlement sync fallback failed: %v", fallbackErr)
	}
	c.JSON(http.StatusOK, res)
}

// GetReceivableByID gets receivable detail.
func GetReceivableByID(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "receivable id required"})
		return
	}

	r, err := SupplyInstance.GetReceivable(nil, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if r.ReceivableId == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "receivable not found"})
		return
	}

	recs, err := SupplyInstance.GetFinancingRecords(nil, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	financing := make([]financingView, 0, len(recs))
	for _, x := range recs {
		financing = append(financing, financingToView(x))
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      receivableToView(r),
		"financing": financing,
	})
}

// GetFinancingByReceivable gets financing list.
func GetFinancingByReceivable(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "receivable id required"})
		return
	}
	recs, err := SupplyInstance.GetFinancingRecords(nil, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	list := make([]financingView, 0, len(recs))
	for _, x := range recs {
		list = append(list, financingToView(x))
	}
	c.JSON(http.StatusOK, gin.H{"count": len(list), "data": list})
}

// ListReceivables lists all receivables from chain.
func ListReceivables(c *gin.Context) {
	count, err := SupplyInstance.GetReceivableCount(nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if count.Sign() == 0 {
		c.JSON(http.StatusOK, gin.H{"count": 0, "data": []receivableView{}})
		return
	}

	total := int(count.Int64())
	list := make([]receivableView, 0, total)
	for i := 0; i < total; i++ {
		id, idxErr := SupplyInstance.GetReceivableIdByIndex(nil, big.NewInt(int64(i)))
		if idxErr != nil {
			continue
		}
		r, recErr := SupplyInstance.GetReceivable(nil, id)
		if recErr != nil || r.ReceivableId == "" {
			continue
		}
		list = append(list, receivableToView(r))
	}

	c.JSON(http.StatusOK, gin.H{"count": len(list), "data": list})
}

// GetOverview returns dashboard data.
func GetOverview(c *gin.Context) {
	users, _ := dbs.ListUsers()
	txLogs, _ := dbs.ListRecentTxLogs(30)

	statusCounters := map[string]int{
		"PENDING":   0,
		"CONFIRMED": 0,
		"FINANCED":  0,
		"SETTLED":   0,
	}
	totalAmount := big.NewInt(0)

	count, err := SupplyInstance.GetReceivableCount(nil)
	if err == nil {
		total := int(count.Int64())
		for i := 0; i < total; i++ {
			id, idxErr := SupplyInstance.GetReceivableIdByIndex(nil, big.NewInt(int64(i)))
			if idxErr != nil {
				continue
			}
			r, getErr := SupplyInstance.GetReceivable(nil, id)
			if getErr != nil || r.ReceivableId == "" {
				continue
			}
			totalAmount = totalAmount.Add(totalAmount, r.Amount)
			statusCounters[statusText(r.Status)]++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"users":          len(users),
		"receivableSize": statusCounters["PENDING"] + statusCounters["CONFIRMED"] + statusCounters["FINANCED"] + statusCounters["SETTLED"],
		"statusCounters": statusCounters,
		"totalAmount":    totalAmount.String(),
		"recentTxLogs":   txLogs,
	})
}

// GetRecentTxLogs returns latest tx logs.
func GetRecentTxLogs(c *gin.Context) {
	page := 1
	pageSize := 30
	if raw := strings.TrimSpace(c.Query("page")); raw != "" {
		if x, err := strconv.Atoi(raw); err == nil && x > 0 {
			page = x
		}
	}
	if raw := strings.TrimSpace(c.Query("pageSize")); raw != "" {
		if x, err := strconv.Atoi(raw); err == nil && x > 0 && x <= 100 {
			pageSize = x
		}
	}
	logs, err := dbs.ListTxLogsPage(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	total, err := dbs.CountTxLogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"count":    len(logs),
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
		"data":     logs,
	})
}

func AnalyzeTxHash(c *gin.Context) {
	hashRaw := strings.TrimSpace(c.Param("hash"))
	if hashRaw == "" || !common.IsHexHash(hashRaw) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tx hash"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	txHash := common.HexToHash(hashRaw)
	txDetail, txErr := SupplyClient.GetTransactionByHash(ctx, txHash)
	receipt, receiptErr := SupplyClient.GetTransactionReceipt(ctx, txHash)

	if txErr != nil && receiptErr != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	result := gin.H{
		"hash":    hashRaw,
		"summary": gin.H{},
	}
	txLog, logErr := dbs.GetTxLogByHash(hashRaw)
	if logErr == nil && txLog.ID != 0 {
		result["audit"] = gin.H{
			"id":           txLog.ID,
			"action":       txLog.Action,
			"businessId":   txLog.BusinessID,
			"status":       txLog.Status,
			"message":      txLog.Message,
			"operatorAddr": txLog.OperatorAddr,
			"createdAt":    txLog.CreatedAt,
		}
	}
	if txDetail != nil {
		result["transaction"] = gin.H{
			"blockHash":        txDetail.GetBlockHash(),
			"blockNumber":      txDetail.GetBlockNumber(),
			"from":             txDetail.GetFrom(),
			"to":               txDetail.GetTo(),
			"gas":              txDetail.GetGas(),
			"gasPrice":         txDetail.GetGasPrice(),
			"nonce":            txDetail.GetNonce(),
			"transactionIndex": txDetail.GetTransactionIndex(),
			"value":            txDetail.GetValue(),
			"inputPreview":     truncateHex(txDetail.GetInput(), 66),
		}
	}
	if receipt != nil {
		logsSummary := make([]gin.H, 0, len(receipt.Logs))
		for index, item := range receipt.Logs {
			if index >= 5 {
				break
			}
			logsSummary = append(logsSummary, gin.H{
				"address":     item.Address,
				"topics":      item.Topics,
				"topicCount":  len(item.Topics),
				"dataPreview": truncateHex(item.Data, 66),
			})
		}
		result["receipt"] = gin.H{
			"blockNumber":     receipt.GetBlockNumber(),
			"gasUsed":         receipt.GetGasUsed(),
			"status":          receipt.GetStatus(),
			"from":            receipt.GetFrom(),
			"to":              receipt.GetTo(),
			"contractAddress": receipt.GetContractAddress().Hex(),
			"outputPreview":   truncateHex(receipt.GetOutput(), 66),
			"logCount":        len(receipt.Logs),
			"logsSummary":     logsSummary,
		}
		result["summary"] = gin.H{
			"statusText": func() string {
				if receipt.GetStatus() == 0 {
					return "执行成功"
				}
				return fmt.Sprintf("执行异常(%d)", receipt.GetStatus())
			}(),
			"blockNumber": receipt.GetBlockNumber(),
			"from":        receipt.GetFrom(),
			"to":          receipt.GetTo(),
			"gasUsed":     receipt.GetGasUsed(),
			"logCount":    len(receipt.Logs),
			"businessId": func() string {
				if txLog.ID != 0 {
					return txLog.BusinessID
				}
				return ""
			}(),
		}
	}
	c.JSON(http.StatusOK, result)
}

func truncateHex(value string, max int) string {
	raw := strings.TrimSpace(value)
	if len(raw) <= max {
		return raw
	}
	if max < 10 {
		return raw[:max]
	}
	return raw[:max] + "..."
}

func parseNodeIDs(raw []byte) []string {
	var ids []string
	if err := json.Unmarshal(raw, &ids); err == nil {
		return ids
	}

	var generic []interface{}
	if err := json.Unmarshal(raw, &generic); err != nil {
		return []string{}
	}
	ids = make([]string, 0, len(generic))
	for _, v := range generic {
		if x := strings.TrimSpace(fmt.Sprintf("%v", v)); x != "" {
			ids = append(ids, x)
		}
	}
	return ids
}

// GetChainStats returns blockchain runtime metrics for admin.
func GetChainStats(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	out := gin.H{
		"groupId":       "",
		"chainId":       "",
		"fiscoVersion":  "",
		"blockHeight":   int64(0),
		"totalTxCount":  "0",
		"failedTxCount": "0",
		"nodeCount":     0,
		"nodeIds":       []string{},
		"warnings":      []string{},
	}

	warnings := make([]string, 0, 4)
	if gid := SupplyClient.GetGroupID(); gid != nil {
		out["groupId"] = gid.String()
	}

	blockHeight, err := SupplyClient.GetBlockNumber(ctx)
	if err != nil {
		warnings = append(warnings, "get block height failed: "+err.Error())
	} else {
		out["blockHeight"] = blockHeight
	}

	txCount, err := SupplyClient.GetTotalTransactionCount(ctx)
	if err != nil {
		warnings = append(warnings, "get tx count failed: "+err.Error())
	} else {
		out["totalTxCount"] = txCount.GetTxSum()
		out["failedTxCount"] = txCount.GetFailedTxSum()
		if x := strings.TrimSpace(txCount.GetBlockNumber()); x != "" {
			if parsed, parseErr := strconv.ParseInt(x, 10, 64); parseErr == nil {
				out["blockHeight"] = parsed
			}
		}
	}

	nodeRaw, err := SupplyClient.GetNodeIDList(ctx)
	if err != nil {
		warnings = append(warnings, "get node ids failed: "+err.Error())
	} else {
		nodeIDs := parseNodeIDs(nodeRaw)
		out["nodeIds"] = nodeIDs
		out["nodeCount"] = len(nodeIDs)
	}

	version, err := SupplyClient.GetClientVersion(ctx)
	if err != nil {
		warnings = append(warnings, "get client version failed: "+err.Error())
	} else {
		out["chainId"] = strings.TrimSpace(version.GetChainId())
		out["fiscoVersion"] = strings.TrimSpace(version.GetFiscoBcosVersion())
	}

	out["warnings"] = warnings
	c.JSON(http.StatusOK, out)
}
