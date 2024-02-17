package telegram

import (
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
	dispatcher.AddHandler(&LinkHandler{bot: bot, svc: service.NewLinkResolutionService(cl)})
	updater := gobot.NewUpdater(dispatcher, nil)
	err = updater.AddWebhook(bot, c.UrlPath, &gobot.AddWebhookOpts{SecretToken: c.SharedSecret})
	if err != nil {
		return nil, err
	}

	return &WebhookBot{b: bot, dispatcher: dispatcher, updater: updater}, nil
}

func (b *WebhookBot) Stop() error {

	return b.updater.Stop()
}

type LinkHandler struct {
	bot *gotgbot.Bot
	svc *service.LinkResolutionService
}

// CheckUpdate checks whether the update should handled by this handler.
func (h *LinkHandler) CheckUpdate(b *gotgbot.Bot, ctx *gobot.Context) bool {
	senderIsNotBot := ctx.EffectiveUser.Id != b.Id
	return senderIsNotBot
}

// HandleUpdate processes the update.
func (h *LinkHandler) HandleUpdate(b *gotgbot.Bot, ctx *gobot.Context) error {
	message := ctx.EffectiveMessage.Text
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
