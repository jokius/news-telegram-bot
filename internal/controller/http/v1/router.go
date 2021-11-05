// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jokius/news-telegram-bot/internal/usecase"
	"github.com/jokius/news-telegram-bot/pkg/logger"
)

func NewRouter(handler *gin.Engine, l logger.InterfaceLogger, user usecase.User, token string) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Routers
	h := handler.Group("/v1")
	{
		UserTelegramRoutes(h, user, token, l)
	}
}
