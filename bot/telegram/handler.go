package telegram

import (
	"errors"
	"github.com/PaulSonOfLars/gotgbot/v2"
	gobot "github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/henges/trackrouter/format"
	"github.com/henges/trackrouter/service"
	"github.com/henges/trackrouter/service/helpers"
	"github.com/rs/zerolog/log"
)

type LinkHandler struct {
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
	user := ctx.EffectiveSender.User.Username

	log.Trace().
		Str("messageBody", message).
		Str("username", user).
		Msg("Handle update")
	result, err := h.svc.FindLinks(message)
	if err != nil {
		// Not an error case.
		if errors.Is(err, helpers.ErrNoMatch) {
			log.Trace().
				Str("messageBody", message).
				Str("username", user).
				Msg("Message didn't match regex")
			return nil
		}
		return err
	}
	if result.Links.IsEmpty() {
		// No matches. Not an error case.
		log.Trace().
			Str("messageBody", message).
			Str("username", user).
			Msg("No matches found")
		return nil
	}
	log.Info().
		Stringer("providerType", result.Id.ProviderType).
		Int("matches", result.Links.Count()).
		Str("username", user).
		Msg("Handled update")
	_, err = b.SendMessage(ctx.EffectiveChat.Id, format.LinksMatchResult(result), nil)
	if err != nil {
		return err
	}
	return nil
}

// Name gets the handler name; used to differentiate handlers programmatically. Names should be unique.
func (h *LinkHandler) Name() string {
	return "LinkHandler"
}
