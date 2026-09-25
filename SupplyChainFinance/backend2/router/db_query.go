package router

import (
	"backend/dbs"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func parseLimit(c *gin.Context, def int) int {
	limit := def
	if raw := strings.TrimSpace(c.Query("limit")); raw != "" {
		if x, err := strconv.Atoi(raw); err == nil && x > 0 {
			limit = x
		}
	}
	return limit
}

// GetDBOverview returns dashboard data from database tables.
func GetDBOverview(c *gin.Context) {
	data, err := dbs.GetDBOverview()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// ListDBEnterprises returns enterprise list from database.
func ListDBEnterprises(c *gin.Context) {
	items, err := dbs.ListDBEnterprises(parseLimit(c, 500))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": len(items), "data": items})
}

// ListDBReceivables returns receivable list from database.
func ListDBReceivables(c *gin.Context) {
	items, err := dbs.ListDBReceivables(parseLimit(c, 1000))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": len(items), "data": items})
}

// GetDBReceivableByID returns receivable detail and financing records from database.
func GetDBReceivableByID(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "receivable id required"})
		return
	}
	rec, err := dbs.GetDBReceivable(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rec.ReceivableID == "" {
		c.JSON(http.StatusOK, gin.H{"data": nil, "financing": []any{}, "message": "receivable not found in database"})
		return
	}
	financing, err := dbs.ListDBFinancingByReceivable(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rec, "financing": financing})
}

// ListDBFinancingByReceivable returns financing list from database.
func ListDBFinancingByReceivable(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "receivable id required"})
		return
	}
	items, err := dbs.ListDBFinancingByReceivable(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": len(items), "data": items})
}

// ListDBDocuments returns document metadata list from database.
func ListDBDocuments(c *gin.Context) {
	items, err := dbs.ListDBDocuments(parseLimit(c, 200))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": len(items), "data": items})
}
