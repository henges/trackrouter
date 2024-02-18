package plexprovider

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/henges/trackrouter/model"
	"github.com/henges/trackrouter/providers"
	"github.com/henges/trackrouter/util"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func NewPlexProvider() providers.Provider {

	return &PlexProvider{&http.Client{}}
}

type PlexProvider struct {
	client *http.Client
}

var plexUrl = "https://listen.plex.tv/track/"
var plexRegex = regexp.MustCompile("listen\\.plex\\.tv/track/(\\w+)")

// MatchId parses text and tries to locate an ID for a track
// that originates from this provider. It should return a wrapped
// ErrMessageNotMatched when no match is found.
func (p *PlexProvider) MatchId(text string) (model.ExternalTrackId, error) {

	if match := util.RegexpMatchWithGroup(text, plexRegex); match != "" {
		return model.ExternalTrackId{
			ProviderType: model.ProviderTypePlex,
			Id:           match,
		}, nil
	}
	return providers.DefaultNoMatchResult(text)
}

// LookupId attempts to locate track metadata based on an ID.
func (p *PlexProvider) LookupId(id string) (model.TrackMetadata, error) {

	URL, err := url.Parse(plexUrl + id)
	if err != nil {
		return model.TrackMetadata{}, err
	}
	req := &http.Request{
		Method: "GET",
		URL:    URL,
	}
	res, err := p.client.Do(req)
	if err != nil {
		return model.TrackMetadata{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return model.TrackMetadata{}, fmt.Errorf("bad response from Plex: %d", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return model.TrackMetadata{}, err
	}
	return parseBody(string(body))
}

func parseBody(body string) (model.TrackMetadata, error) {

	r := strings.NewReader(body)
	reader, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return model.TrackMetadata{}, err
	}

	trackTitle := reader.Find(".titles h1").Text()
	subtitles := reader.Find(".titles h2 a").Nodes
	if len(subtitles) < 2 {
		return model.TrackMetadata{}, errors.New("missing html element")
	}
	album := subtitles[0].FirstChild.Data
	artist := subtitles[1].FirstChild.Data
	return model.TrackMetadata{
		Title:   trackTitle,
		Artists: []string{artist},
		Album:   album,
	}, nil
}

// LookupMetadata attempts to locate a matching ID for a track
// based on the input metadata.
func (p *PlexProvider) LookupMetadata(metadata model.TrackMetadata) string {

	return ""
}
