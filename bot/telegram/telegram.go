package telegram

import (
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	gobot "github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/henges/trackrouter/config"
	"github.com/henges/trackrouter/service"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

type WebhookBot struct {
	b          *gotgbot.Bot
	dispatcher *gobot.Dispatcher
	updater    *gobot.Updater
	c          *config.TelegramConfig
	cmds       Commands
}

type Command struct {
	gotgbot.BotCommand
	Func handlers.Response
}

type Commands []Command

func GetCommands(svc *service.LinkResolutionService) Commands {

	lr := &LinkResponse{svc: svc}

	return Commands{
		{
			BotCommand: gotgbot.BotCommand{
				Command:     "link",
				Description: "Attempts to find links to tracks on streaming services based on a query.",
			},
			Func: lr.Response,
		},
		{
			BotCommand: gotgbot.BotCommand{
				Command:     "help",
				Description: "Displays information about this bot.",
			},
			Func: func(b *gotgbot.Bot, ctx *gobot.Context) error {

				_, err := ctx.EffectiveChat.SendMessage(b, "trackrouter helps find links to tracks on streaming services.", nil)
				return err
			},
		},
	}
}

func CommandsEqual(v1 []Command, v2 []gotgbot.BotCommand) bool {

	if len(v1) != len(v2) {
		return false
	}
	for i, v := range v1 {
		if v2[i] != v.BotCommand {
			return false
		}
	}

	return true
}

func (c Commands) GetGobotCommands() []gotgbot.BotCommand {

	return lo.Map(c, func(item Command, index int) gotgbot.BotCommand {
		return item.BotCommand
	})
}

func NewWebhookBot(c *config.TelegramConfig, svc *service.LinkResolutionService) (*WebhookBot, error) {

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

	cmds := GetCommands(svc)
	for _, v := range cmds {
		dispatcher.AddHandler(handlers.NewCommand(v.Command, v.Func))
	}
	dispatcher.AddHandler(&URLHandler{svc: svc})
	updater := gobot.NewUpdater(dispatcher, nil)
	err = updater.AddWebhook(bot, c.UrlPath, &gobot.AddWebhookOpts{SecretToken: c.SharedSecret})
	if err != nil {
		return nil, err
	}

	return &WebhookBot{b: bot, dispatcher: dispatcher, updater: updater, c: c, cmds: cmds}, nil
}

func (b *WebhookBot) Start() error {
	oldCommands, err := b.b.GetMyCommands(nil)
	if err != nil {
		return err
	}
	if !CommandsEqual(b.cmds, oldCommands) {

		ok, err := b.b.SetMyCommands(b.cmds.GetGobotCommands(), nil)
		if err != nil {
			return err
		}
		if !ok {
			log.Error().Msg("Non ok result when trying to update commands")
			return nil
		}
		log.Info().Msg("Updated commands")
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
