package httpapi

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type ReadinessCheck func(context.Context) error

func init() {
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}
}

func Handler(check ReadinessCheck, logger *slog.Logger) http.Handler {
	router := gin.New()
	router.Use(gin.Recovery())
	registerHealth(router, check, logger)
	return router
}

func registerHealth(router *gin.Engine, check ReadinessCheck, logger *slog.Logger) {
	router.GET("/healthz", func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.JSON(200, gin.H{"service": "api_gateway", "status": "ok"})
	})
	router.GET("/readyz", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := check(ctx); err != nil {
			logger.Warn("readiness check failed", "reason", err.Error())
			c.Header("Cache-Control", "no-store")
			c.JSON(503, gin.H{"service": "api_gateway", "status": "not_ready"})
			return
		}
		c.Header("Cache-Control", "no-store")
		c.JSON(200, gin.H{"service": "api_gateway", "status": "ready"})
	})
}
