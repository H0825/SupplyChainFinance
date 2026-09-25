package router

import (
	"backend/SupplyChainFinance"
	"backend/dbs"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/FISCO-BCOS/go-sdk/client"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var ks *keystore.KeyStore

func normalizeRole(role string) string {
	r := strings.ToLower(strings.TrimSpace(role))
	switch r {
	case "business", "finance", "admin":
		return r
	default:
		return ""
	}
}

func init() {
	keyDir := "./keystore"
	if err := os.MkdirAll(keyDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create keystore directory: %v", err)
	}
	ks = keystore.NewKeyStore(keyDir, keystore.LightScryptN, keystore.LightScryptP)
}

// UserRegister handles user registration.
func UserRegister(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if strings.TrimSpace(req.Role) == "" {
		req.Role = "business"
	}
	req.Role = normalizeRole(req.Role)
	if req.Role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role must be one of: business, finance, admin"})
		return
	}

	exists, err := dbs.GetUser(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query user"})
		return
	}
	if exists.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
		return
	}

	account, err := ks.NewAccount(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create account"})
		return
	}

	keyJSON, err := ks.Export(account, req.Password, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to export account"})
		return
	}
	decrypted, err := keystore.DecryptKey(keyJSON, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decrypt key"})
		return
	}

	if err = dbs.SaveUser(req.Username, req.Password, account.Address.Hex(), req.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "register successful",
		"address":    account.Address.Hex(),
		"privateKey": hex.EncodeToString(decrypted.PrivateKey.D.Bytes()),
	})
}

// UserLogin unlocks account and binds transactor with current key.
func UserLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := dbs.GetUser(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query user"})
		return
	}
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	if !common.IsHexAddress(user.Address) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid wallet address in db"})
		return
	}
	addr := common.HexToAddress(user.Address)

	acc, err := ks.Find(accounts.Account{Address: addr})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find account in local keystore"})
		return
	}

	if err = ks.TimedUnlock(acc, req.Password, 5*time.Minute); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to unlock account"})
		return
	}

	keyJSON, err := ks.Export(acc, req.Password, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to export account"})
		return
	}
	decrypted, err := keystore.DecryptKey(keyJSON, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decrypt account key"})
		return
	}

	Configs[0].PrivateKey = decrypted.PrivateKey.D.Bytes()
	SupplyClient, err = client.Dial(&Configs[0])
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect blockchain client"})
		return
	}
	SupplyInstance, err = SupplyChainFinance.NewSupplyChainFinance(ContractAddress, SupplyClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create contract instance"})
		return
	}

	fmt.Println("user", addr.Hex(), "login successful")
	c.JSON(http.StatusOK, gin.H{
		"message":  "login successful",
		"address":  addr.Hex(),
		"role":     user.Role,
		"user_id":  user.ID,
		"username": user.Username,
	})
}

// ListUsers provides user list for front-end dropdown and mapping.
func ListUsers(c *gin.Context) {
	users, err := dbs.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": len(users), "data": users})
}
