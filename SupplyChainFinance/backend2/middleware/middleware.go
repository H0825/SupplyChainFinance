package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type routeStat struct {
	TotalRequests int     `json:"totalRequests"`
	SuccessCount  int     `json:"successCount"`
	FailureCount  int     `json:"failureCount"`
	TotalTime     int64   `json:"totalTime"`
	AverageTime   float64 `json:"averageTime"`
}

var (
	statsMu       sync.Mutex
	routeStats    = make(map[string]routeStat)
	totalRequests int
)

// CORSMiddleware handles cross-origin requests.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-Role")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// StatsMiddleware captures route-level metrics.
func StatsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}
		start := time.Now()

		c.Next()

		elapsed := time.Since(start).Milliseconds()
		status := c.Writer.Status()

		statsMu.Lock()
		defer statsMu.Unlock()

		totalRequests++
		st := routeStats[route]
		st.TotalRequests++
		st.TotalTime += elapsed
		st.AverageTime = float64(st.TotalTime) / float64(st.TotalRequests)
		if status >= 200 && status < 400 {
			st.SuccessCount++
		} else {
			st.FailureCount++
		}
		routeStats[route] = st

		log.Printf("route=%s status=%d cost_ms=%d total=%d", route, status, elapsed, totalRequests)
	}
}

// GetStats returns middleware stats.
func GetStats(c *gin.Context) {
	statsMu.Lock()
	defer statsMu.Unlock()

	snapshot := make(map[string]routeStat, len(routeStats))
	for k, v := range routeStats {
		snapshot[k] = v
	}

	c.JSON(http.StatusOK, gin.H{
		"totalRequests": totalRequests,
		"routeStats":    snapshot,
		"timestamp":     time.Now().Format("2006-01-02 15:04:05"),
	})
}

// Recovery catches panics and returns 500.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
				c.Abort()
			}
		}()
		c.Next()
	}
}
