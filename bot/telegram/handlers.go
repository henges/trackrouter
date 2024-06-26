package telegram

import (
	"errors"
	"github.com/PaulSonOfLars/gotgbot/v2"
	gobot "github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/henges/trackrouter/format"
	"github.com/henges/trackrouter/providers"
	"github.com/henges/trackrouter/service"
	"github.com/rs/zerolog/log"
	"strings"
)

type URLHandler struct {
	svc *service.LinkResolutionService
}

// CheckUpdate checks whether the update should handled by this handler.
func (h *URLHandler) CheckUpdate(b *gotgbot.Bot, ctx *gobot.Context) bool {
	senderIsNotBot := ctx.EffectiveUser.Id != b.Id
	return senderIsNotBot
}

// HandleUpdate processes the update.
func (h *URLHandler) HandleUpdate(b *gotgbot.Bot, ctx *gobot.Context) error {
	message := ctx.EffectiveMessage.Text
	user := ctx.EffectiveSender.User.Username

	log.Trace().
		Str("messageBody", message).
		Str("username", user).
		Msg("Handle update")
	result, err := h.svc.FindLinksFromUrl(message)
	if err != nil {
		// Not an error case.
		if errors.Is(err, providers.ErrMessageNotMatched) {
			log.Trace().
				Err(err).
				Str("messageBody", message).
				Str("username", user).
				Msg("Message didn't match regex")
			return nil
		}
		if errors.Is(err, providers.ErrIdNotMatched) {
			log.Trace().
				Err(err).
				Str("messageBody", message).
				Str("username", user).
				Msg("Couldn't find any matches for track")
			return nil
		}
		return err
	}
	log.Info().
		Stringer("providerType", result.Id.ProviderType).
		Any("metadata", result.TrackMetadata).
		Int("matches", len(result.Links)).
		Str("username", user).
		Msg("Handled update")
	_, err = b.SendMessage(ctx.EffectiveChat.Id, format.LinksMatchResult(result.Links), nil)
	if err != nil {
		return err
	}
	return nil
}

// Name gets the handler name; used to differentiate handlers programmatically. Names should be unique.
func (h *URLHandler) Name() string {
	return "URLHandler"
}

type LinkResponse struct {
	svc    *service.LinkResolutionService
	prefix string
}

func (h *LinkResponse) Response(b *gotgbot.Bot, ctx *gobot.Context) error {
	message := ctx.EffectiveMessage.Text
	user := ctx.EffectiveSender.User.Username

	log.Trace().
		Str("messageBody", message).
		Str("username", user).
		Msg("Handle update")

	// Strip command name
	message = strings.TrimSpace(strings.Replace(message, "/"+h.prefix, "", 1))
	result, err := h.svc.FindLinksFromMessage(message)
	if err != nil {
		// Not an error case.
		if errors.Is(err, providers.ErrMessageNotMatched) {
			log.Trace().
				Err(err).
				Str("messageBody", message).
				Str("username", user).
				Msg("Message didn't match regex")
			return nil
		}
		if errors.Is(err, providers.ErrIdNotMatched) {
			log.Trace().
				Err(err).
				Str("messageBody", message).
				Str("username", user).
				Msg("Couldn't find any matches for track")
			return nil
		}
		return err
	}
	log.Info().
		Str("query", message).
		Int("matches", len(result)).
		Str("username", user).
		Msg("Handled update")
	_, err = b.SendMessage(ctx.EffectiveChat.Id, format.LinksMatchResult(result), nil)
	if err != nil {
		return err
	}
	return nil
}
