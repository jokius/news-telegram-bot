package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jokius/news-telegram-bot/internal/entity"
	"github.com/jokius/news-telegram-bot/internal/usecase"
	"github.com/jokius/news-telegram-bot/pkg/logger"
)

type userRoutes struct {
	user usecase.User
	l    logger.InterfaceLogger
}

func UserTelegramRoutes(handler *gin.RouterGroup, u usecase.User, token string, l logger.InterfaceLogger) {
	r := &userRoutes{u, l}

	h := handler.Group("/telegram")
	{
		h.POST("/callback/"+token, r.telegramCallback)
	}
}

func (r *userRoutes) telegramCallback(c *gin.Context) {
	var telegramResult entity.TelegramResult
	err := c.ShouldBindJSON(&telegramResult)

	if err == nil {
		err = r.user.TelegramCallback(telegramResult)
		if err != nil {
			r.l.Error(fmt.Errorf("`r.telegramCallback` something wrong: %w", err))
		}
	} else {
		r.l.Error(fmt.Errorf("incorrect json: %w", err))
	}

	c.Status(http.StatusNoContent)
}
