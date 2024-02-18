package format

import (
	"github.com/henges/trackrouter/model"
	"strings"
)

func LinksMatchResult(l *model.LinksMatchResult) string {

	// The order we want links to appear in
	order := []model.ProviderType{
		model.ProviderTypeYoutube, // Youtube first since it provides the link preview
		model.ProviderTypeTidal,
		model.ProviderTypeSpotify,
		model.ProviderTypePlex,
	}
	strs := make([]string, 0, len(l.Links))
	for _, v := range order {
		if val, ok := l.Links[v]; ok {
			strs = append(strs, val)
		}
	}

	return strings.Join(strs, "\n")
}
