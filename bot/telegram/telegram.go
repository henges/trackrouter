package telegram

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	gobot "github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/henges/trackrouter/config"
	"github.com/henges/trackrouter/service"
	"github.com/rs/zerolog/log"
)

type WebhookBot struct {
	b          *gotgbot.Bot
	dispatcher *gobot.Dispatcher
	updater    *gobot.Updater
	c          *config.TelegramConfig
}

func NewWebhookBot(c *config.TelegramConfig, lres *service.LinkResolutionService) (*WebhookBot, error) {

	bot, err := gotgbot.NewBot(c.AuthToken, nil)
	if err != nil {
		return nil, err
	}

	dispatcher := gobot.NewDispatcher(&gobot.DispatcherOpts{
		// If an error is returned by a handler, log it and continue going.
		Error: func(b *gotgbot.Bot, ctx *gobot.Context, err error) gobot.DispatcherAction {
			log.Err(err).Msg("error occurred while handling update")
			return gobot.DispatcherActionNoop
		},
		MaxRoutines: gobot.DefaultMaxRoutines,
	})
	dispatcher.AddHandler(&LinkHandler{svc: lres})
	updater := gobot.NewUpdater(dispatcher, nil)
	err = updater.AddWebhook(bot, c.UrlPath, &gobot.AddWebhookOpts{SecretToken: c.SharedSecret})
	if err != nil {
		return nil, err
	}

	return &WebhookBot{b: bot, dispatcher: dispatcher, updater: updater, c: c}, nil
}

func (b *WebhookBot) Start() error {
	err := b.updater.StartServer(gobot.WebhookOpts{ListenAddr: fmt.Sprintf("0.0.0.0:%d", b.c.ListenPort), SecretToken: b.c.SharedSecret})
	if err != nil {
		return err
	}
	return b.updater.SetAllBotWebhooks(b.c.Host, &gotgbot.SetWebhookOpts{SecretToken: b.c.SharedSecret})
}

func (b *WebhookBot) Stop() error {

	return b.updater.Stop()
}
