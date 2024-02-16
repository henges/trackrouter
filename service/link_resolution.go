package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/henges/trackrouter/di"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/sync/errgroup"
	"regexp"
)

var spotifyRegex = regexp.MustCompile("open\\.spotify\\.com/track/(\\w+)")
var tidalRegex = regexp.MustCompile("tidal\\.com/track/(\\d+)")

type LinkResolutionService struct {
	Clients *di.Clients
}

func NewLinkResolutionService(di *di.Deps) *LinkResolutionService {

	return &LinkResolutionService{Clients: di.Clients}
}

type Links struct {
	SpotifyLink string `json:"spotifyLink"`
	TidalLink   string `json:"tidalLink"`
	YoutubeLink string `json:"youtubeLink"`
}

func (l *LinkResolutionService) GetLinks(query string) (*Links, error) {

	var eg errgroup.Group
	var links Links

	eg.Go(func() error {
		search, err := l.Clients.SpotifyClient.Search(context.TODO(), query, spotify.SearchTypeTrack)
		if err != nil {
			return err
		}
		if len(search.Tracks.Tracks) > 0 {
			links.SpotifyLink = fmt.Sprintf("https://open.spotify.com/track/%s", search.Tracks.Tracks[0].ID)
		}
		return nil
	})
	eg.Go(func() error {
		search, err := l.Clients.TidalClient.Search(context.TODO(), query)
		if err != nil {
			return err
		}
		if search.Tracks == nil {
			return errors.New("tracks was nil in tidal response")
		}
		tracks := *search.Tracks
		if len(tracks) > 0 {
			links.TidalLink = tracks[0].Resource.TidalUrl
		}
		return nil
	})

	err := eg.Wait()
	return &links, err
}

type ProviderType int

const (
	ProviderTypeSpotify ProviderType = 1
	ProviderTypeTidal   ProviderType = 2
)

type ExternalTrackId struct {
	ProviderType ProviderType
	Id           string
}

var ErrNoMatch = errors.New("no match found for input text")

func ResolveId(text string) (ExternalTrackId, error) {

	spotifyMatch := regexpMatchWithGroup(text, spotifyRegex)
	if spotifyMatch != "" {
		return ExternalTrackId{ProviderType: ProviderTypeSpotify, Id: spotifyMatch}, nil
	}
	tidalMatch := regexpMatchWithGroup(text, tidalRegex)
	if tidalMatch != "" {
		return ExternalTrackId{ProviderType: ProviderTypeTidal, Id: tidalMatch}, nil
	}

	return ExternalTrackId{}, fmt.Errorf("for input text %s: %w", text, ErrNoMatch)
}

func regexpMatchWithGroup(text string, exp *regexp.Regexp) string {

	matches := exp.FindStringSubmatch(text)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}
