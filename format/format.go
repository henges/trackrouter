package format

import (
	"github.com/henges/trackrouter/model"
	"strings"
)

func LinksMatchResult(l model.Links) string {

	// The order we want links to appear in
	order := []model.ProviderType{
		model.ProviderTypeYoutube, // Youtube first since it provides the link preview
		model.ProviderTypeTidal,
		model.ProviderTypeSpotify,
		model.ProviderTypePlex,
	}
	strs := make([]string, 0, len(l))
	for _, v := range order {
		if val, ok := l[v]; ok {
			strs = append(strs, val)
		}
	}

	return strings.Join(strs, "\n")
}
