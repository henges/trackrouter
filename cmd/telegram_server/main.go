package main

import (
	"context"
	"github.com/henges/trackrouter/bot/telegram"
	"github.com/henges/trackrouter/config"
	"github.com/henges/trackrouter/di"
	"github.com/rs/zerolog/log"
	"os/signal"
	"syscall"
)

func main() {

	c := config.Get()
	log.Info().Msg("App started")

	deps := di.Get(c)
	b, err := telegram.NewWebhookBot(c.Telegram, deps.Clients)
	if err != nil {
		log.Fatal().Err(err).Msg("while creating telegram bot")
	}

	err = b.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("while starting telegram bot")
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	log.Info().Msg("App ready")

	<-ctx.Done()
	stop()
	err = b.Stop()
	log.Info().Err(err).Msg("App shutdown")
}
