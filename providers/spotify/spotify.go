package spotifyprovider

import (
	"context"
	"fmt"
	"github.com/henges/trackrouter/model"
	"github.com/henges/trackrouter/providers"
	"github.com/henges/trackrouter/util"
	"github.com/samber/lo"
	"github.com/zmb3/spotify/v2"
	"regexp"
)

func NewSpotifyProvider(c *spotify.Client) providers.Provider {

	return &SpotifyProvider{
		SpotifyMatch:  &SpotifyMatch{},
		SpotifyLookup: &SpotifyLookup{c},
	}
}

type SpotifyProvider struct {
	*SpotifyMatch
	*SpotifyLookup
}

type SpotifyMatch struct{}

type SpotifyLookup struct {
	client *spotify.Client
}

var spotifyRegex = regexp.MustCompile("open\\.spotify\\.com/track/(\\w+)")

func (s *SpotifyMatch) MatchId(text string) (model.ExternalTrackId, error) {

	if match := util.RegexpMatchWithGroup(text, spotifyRegex); match != "" {
		return model.ExternalTrackId{ProviderType: model.ProviderTypeSpotify, Id: match}, nil
	}

	return providers.DefaultNoMatchResult(text)
}

func (s *SpotifyLookup) LookupId(id string) (model.TrackMetadata, error) {

	track, err := s.client.GetTrack(context.TODO(), spotify.ID(id))
	if err != nil {
		return model.TrackMetadata{}, err
	}
	return model.TrackMetadata{
		Title: track.Name,
		Artists: lo.Map(track.Artists, func(item spotify.SimpleArtist, index int) string {
			return item.Name
		}),
		Album: track.Album.Name,
	}, nil
}

func (s *SpotifyLookup) LookupMetadata(metadata model.TrackMetadata) string {

	if result, err := s.doSearch(providers.DefaultTrackMetadataQuery(metadata)); err != nil {
		return ""
	} else {
		return result
	}
}

func (s *SpotifyLookup) LookupQuery(query string) string {

	if result, err := s.doSearch(query); err != nil {
		return ""
	} else {
		return result
	}
}

func (s *SpotifyLookup) doSearch(query string) (string, error) {

	search, err := s.client.Search(context.Background(), query, spotify.SearchTypeTrack)
	if err != nil {
		return "", err
	}
	if len(search.Tracks.Tracks) > 0 {
		return fmt.Sprintf("https://open.spotify.com/track/%s", search.Tracks.Tracks[0].ID), nil
	}
	return "", providers.ErrMessageNotMatched
}
