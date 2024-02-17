package format

import (
	"github.com/henges/trackrouter/model"
	"github.com/samber/lo"
	"strings"
)

func LinksMatchResult(l *model.LinksMatchResult) string {

	var values []string
	switch l.Id.ProviderType {
	case model.ProviderTypeSpotify:
		values = []string{l.Links.YoutubeLink, l.Links.TidalLink}
	case model.ProviderTypeTidal:
		values = []string{l.Links.YoutubeLink, l.Links.SpotifyLink}
	}

	nonEmpty := lo.Filter(values, func(s string, _ int) bool {
		return s != ""
	})
	return strings.Join(nonEmpty, "\n")
}
