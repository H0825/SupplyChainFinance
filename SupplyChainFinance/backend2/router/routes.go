package router

import (
	"backend/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	r.GET("/workflow/public/stats", middleware.GetStats)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", UserRegister)
			auth.POST("/login", UserLogin)
			auth.GET("/users", ListUsers)
		}

		enterprise := api.Group("/enterprise")
		{
			enterprise.POST("/register", middleware.RequireRoles("business", "admin"), RegisterEnterprise)
			enterprise.GET("", ListEnterprises)
			enterprise.GET("/:address", GetEnterpriseByAddress)
		}

		receivable := api.Group("/receivable")
		{
			receivable.POST("/issue", middleware.RequireRoles("business", "admin"), IssueReceivable)
			receivable.POST("/confirm", middleware.RequireRoles("business", "admin"), ConfirmReceivable)
			receivable.POST("/finance", middleware.RequireRoles("finance", "admin"), RecordFinancing)
			receivable.POST("/settle", middleware.RequireRoles("business", "admin"), RecordSettlement)
			receivable.GET("", ListReceivables)
			receivable.GET("/:id", GetReceivableByID)
			receivable.GET("/:id/financing", GetFinancingByReceivable)
		}

		business := api.Group("/business", middleware.RequireRoles("business", "admin"))
		{
			business.POST("/enterprise/register", RegisterEnterprise)
			business.POST("/receivable/issue", IssueReceivable)
			business.POST("/receivable/confirm", ConfirmReceivable)
			business.POST("/receivable/settle", RecordSettlement)
			business.GET("/receivable/:id", GetReceivableByID)
			business.GET("/receivables", ListReceivables)
		}

		finance := api.Group("/finance", middleware.RequireRoles("finance", "admin"))
		{
			finance.POST("/receivable/finance", RecordFinancing)
			finance.GET("/receivable/:id", GetReceivableByID)
			finance.GET("/receivable/:id/financing", GetFinancingByReceivable)
			finance.GET("/receivables", ListReceivables)
		}

		api.GET("/dashboard/overview", GetOverview)
		api.GET("/tx-logs", GetRecentTxLogs)
		api.GET("/tx-logs/:hash/analyze", middleware.RequireRoles("business", "finance", "admin"), AnalyzeTxHash)
		api.GET("/chain/stats", middleware.RequireRoles("business", "finance", "admin"), GetChainStats)
		api.POST("/ai/chat", middleware.RequireRoles("business", "finance", "admin"), ChatAssistant)
		api.POST("/ipfs/upload", middleware.RequireRoles("business", "finance", "admin"), FileUpload)

		db := api.Group("/db")
		{
			db.GET("/dashboard/overview", GetDBOverview)
			db.GET("/enterprises", ListDBEnterprises)
			db.GET("/receivables", ListDBReceivables)
			db.GET("/receivables/:id", GetDBReceivableByID)
			db.GET("/receivables/:id/financing", ListDBFinancingByReceivable)
			db.GET("/documents", ListDBDocuments)
		}
	}
}
