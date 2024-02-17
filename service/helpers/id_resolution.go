package helpers

import (
	"errors"
	"fmt"
	"github.com/henges/trackrouter/model"
	"regexp"
)

var spotifyRegex = regexp.MustCompile("open\\.spotify\\.com/track/(\\w+)")
var tidalRegex = regexp.MustCompile("tidal\\.com/track/(\\d+)")
var ErrNoMatch = errors.New("no match found for input text")

func ResolveId(text string) (model.ExternalTrackId, error) {

	spotifyMatch := regexpMatchWithGroup(text, spotifyRegex)
	if spotifyMatch != "" {
		return model.ExternalTrackId{ProviderType: model.ProviderTypeSpotify, Id: spotifyMatch}, nil
	}
	tidalMatch := regexpMatchWithGroup(text, tidalRegex)
	if tidalMatch != "" {
		return model.ExternalTrackId{ProviderType: model.ProviderTypeTidal, Id: tidalMatch}, nil
	}

	return model.ExternalTrackId{}, fmt.Errorf("for input text %s: %w", text, ErrNoMatch)
}

func regexpMatchWithGroup(text string, exp *regexp.Regexp) string {

	matches := exp.FindStringSubmatch(text)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}
