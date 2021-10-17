package main

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jokius/news-telegram-bot/config"
	"github.com/jokius/news-telegram-bot/internal/app"
)

func main() {
	// Configuration
	var cfg config.Config

	err := cleanenv.ReadConfig("./config/config.yml", &cfg)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(&cfg)
}
