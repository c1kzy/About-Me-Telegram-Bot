package main

import (
	"fmt"
	"net/http"
	"telegram/telegram-api"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"github.com/phuslu/log"
)

func main() {

	if envErr := godotenv.Load(); envErr != nil {
		log.Fatal().Err(envErr).Msgf("error loading .env file", envErr)
	}

	cfg := &telegram_api.Config{}

	log.DefaultLogger = log.Logger{
		Level:      log.InfoLevel,
		Caller:     cfg.Caller,
		TimeField:  cfg.TimeField,
		TimeFormat: time.RFC850,
		Writer:     &log.ConsoleWriter{},
	}

	if err := env.Parse(cfg); err != nil {
		log.Error().Err(err)
	}
	api := telegram_api.GetAPI(cfg)

	http.HandleFunc("/telegram", api.TelegramHandler)

	err := http.ListenAndServe(fmt.Sprintf(":%v", cfg.Port), nil)
	if err != nil {
		log.Fatal().Err(err).Msgf("server start failed %v", err)
	}
	log.Info().Msg("Server started")
}
