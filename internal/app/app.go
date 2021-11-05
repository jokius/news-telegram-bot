// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jokius/news-telegram-bot/config"
	v1 "github.com/jokius/news-telegram-bot/internal/controller/http/v1"
	"github.com/jokius/news-telegram-bot/internal/usecase"
	"github.com/jokius/news-telegram-bot/internal/usecase/repo"
	"github.com/jokius/news-telegram-bot/internal/usecase/service"
	"github.com/jokius/news-telegram-bot/pkg/grabber"
	"github.com/jokius/news-telegram-bot/pkg/httpclient"
	"github.com/jokius/news-telegram-bot/pkg/httpserver"
	"github.com/jokius/news-telegram-bot/pkg/logger"
	"github.com/jokius/news-telegram-bot/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	l.Info(cfg.PG.URL)
	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax),
		postgres.Debug(os.Getenv("DATABASE_DEBUG_INFO")))

	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}

	// Use case
	client := httpclient.NewClient()
	source := service.NewVkSource(cfg.Vk.Token, client)
	messenger := service.NewMessenger(cfg.Telegram.Token, cfg.Telegram.BaseURL, client, source, l)
	userUseCase := usecase.NewUserUseCase(
		repo.NewUserRepo(pg),
		messenger,
		source,
	)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, userUseCase, cfg.Telegram.Token)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Grabbers server
	var apiGrabbers []grabber.Grabber

	sleepTime := time.Duration(cfg.Grabber.Sleep) * time.Second
	groupRepo := repo.NewGroupRepo(pg)
	messageRepo := repo.NewMessageRepo(pg)
	vkGrabber := service.NewVkGrabber(sleepTime, source, messenger, groupRepo, messageRepo, l)
	apiGrabbers = append(apiGrabbers, &vkGrabber)
	grabbersServer := grabber.New(apiGrabbers)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	grabbersServer.Shutdown()
}
