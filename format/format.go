package format

import (
	"github.com/henges/trackrouter/model"
	"github.com/samber/lo"
	"strings"
)

func Links(l *model.Links) string {

	nonEmpty := lo.Filter([]string{l.SpotifyLink, l.TidalLink, l.YoutubeLink}, func(s string, _ int) bool {
		return s != ""
	})

	return strings.Join(nonEmpty, "\n")
}
