package telegram

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	gobot "github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/henges/trackrouter/config"
	"github.com/henges/trackrouter/di"
	"github.com/henges/trackrouter/format"
	"github.com/henges/trackrouter/service"
	"github.com/rs/zerolog/log"
)

type WebhookBot struct {
	b          *gotgbot.Bot
	dispatcher *gobot.Dispatcher
	updater    *gobot.Updater
	c          *config.TelegramConfig
}

func NewWebhookBot(c *config.TelegramConfig, cl *di.Clients) (*WebhookBot, error) {

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
	dispatcher.AddHandler(&LinkHandler{svc: service.NewLinkResolutionService(cl)})
	updater := gobot.NewUpdater(dispatcher, nil)

	return &WebhookBot{b: bot, dispatcher: dispatcher, updater: updater, c: c}, nil
}

func (b *WebhookBot) Start() error {

	err := b.updater.AddWebhook(b.b, b.c.UrlPath, &gobot.AddWebhookOpts{SecretToken: b.c.SharedSecret})
	if err != nil {
		return err
	}
	err = b.updater.StartServer(gobot.WebhookOpts{ListenAddr: fmt.Sprintf("0.0.0.0:%d", b.c.ListenPort), SecretToken: b.c.SharedSecret})
	if err != nil {
		return err
	}
	return b.updater.SetAllBotWebhooks(b.c.Host, &gotgbot.SetWebhookOpts{SecretToken: b.c.SharedSecret})
}

func (b *WebhookBot) Stop() error {

	return b.updater.Stop()
}

type LinkHandler struct {
	svc *service.LinkResolutionService
}

// CheckUpdate checks whether the update should handled by this handler.
func (h *LinkHandler) CheckUpdate(b *gotgbot.Bot, ctx *gobot.Context) bool {
	log.Debug().Msg("Received update")
	senderIsNotBot := ctx.EffectiveUser.Id != b.Id
	return senderIsNotBot
}

// HandleUpdate processes the update.
func (h *LinkHandler) HandleUpdate(b *gotgbot.Bot, ctx *gobot.Context) error {
	message := ctx.EffectiveMessage.Text
	log.Debug().Str("messageBody", message).Str("username", ctx.EffectiveUser.Username).Msg("Handle update")
	links, err := h.svc.FindLinks(message)
	if err != nil {
		return err
	}
	// No matches
	if links.IsEmpty() {
		return nil
	}
	_, err = b.SendMessage(ctx.EffectiveChat.Id, format.Links(links), nil)
	if err != nil {
		return err
	}
	return nil
}

// Name gets the handler name; used to differentiate handlers programmatically. Names should be unique.
func (h *LinkHandler) Name() string {
	return "LinkHandler"
}
