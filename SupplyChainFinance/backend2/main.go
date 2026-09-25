package main

import (
	"backend/appconfig"
	"backend/dbs"
	"backend/middleware"
	"backend/router"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := appconfig.Load("config.toml")
	if err != nil {
		log.Fatal(err)
	}

	dbs.InitDB(cfg.DSN())
	defer dbs.CloseDB()

	f, err := os.OpenFile(cfg.App.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("could not create log file: %v", err)
	}
	defer f.Close()
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.StatsMiddleware())
	r.Use(middleware.Recovery())

	router.RegisterRoutes(r)

	frontendIndex := cfg.App.FrontendDist + "/index.html"
	frontendAssets := cfg.App.FrontendDist + "/assets"
	localIndex := cfg.App.LegacyIndex
	localAssets := cfg.App.LegacyAssets

	staticIndex := localIndex
	staticAssets := localAssets
	if _, statErr := os.Stat(frontendIndex); statErr == nil {
		staticIndex = frontendIndex
		staticAssets = frontendAssets
	}

	r.Static("/assets", staticAssets)
	r.NoRoute(func(c *gin.Context) {
		if c.Request.Method == http.MethodGet {
			c.File(staticIndex)
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	})

	if err = r.Run(cfg.App.Listen); err != nil {
		log.Fatal(err)
	}
}
