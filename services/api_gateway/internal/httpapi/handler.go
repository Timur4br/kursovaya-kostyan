package httpapi

import (
	"context"
	"fmt"
	"io"
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

func Handler(check ReadinessCheck, logger *slog.Logger, dependencies *Controllers) http.Handler {
	router := gin.New()
	_ = router.SetTrustedProxies(nil)
	router.HandleMethodNotAllowed = true
	router.RedirectTrailingSlash = false
	router.Use(gin.CustomRecoveryWithWriter(io.Discard, func(c *gin.Context, recovered any) {
		logger.Error("HTTP panic", "error_type", fmt.Sprintf("%T", recovered))
		c.AbortWithStatusJSON(500, gin.H{"error": gin.H{"code": "internal_error", "message": "Ошибка сервера"}})
	}))
	router.Use(func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})
	registerHealth(router, check, logger)
	router.GET("/openapi.json", func(c *gin.Context) { c.JSON(200, OpenAPI()) })
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": gin.H{"code": "not_found", "message": "Маршрут не найден"}})
	})
	router.NoMethod(func(c *gin.Context) {
		c.JSON(405, gin.H{"error": gin.H{"code": "method_not_allowed", "message": "Метод не поддерживается"}})
	})
	if dependencies != nil {
		registerEntities(router, dependencies, logger)
	}
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
