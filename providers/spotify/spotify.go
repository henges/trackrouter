package spotifyprovider

import (
	"context"
	"fmt"
	"github.com/henges/trackrouter/model"
	"github.com/henges/trackrouter/providers/helpers"
	"github.com/henges/trackrouter/providers/types"
	"github.com/henges/trackrouter/util"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/zmb3/spotify/v2"
	"regexp"
)

func NewSpotifyProvider(c *spotify.Client) providertypes.Provider {

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

	return providerhelpers.DefaultNoMatchResult(text)
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

	query := providerhelpers.DefaultTrackMetadataQuery(metadata)
	search, err := s.client.Search(context.Background(), query, spotify.SearchTypeTrack)
	if err != nil {
		log.Error().Err(err).Msg("in spotify request")
		return ""
	}
	if len(search.Tracks.Tracks) > 0 {
		return fmt.Sprintf("https://open.spotify.com/track/%s", search.Tracks.Tracks[0].ID)
	}
	return ""
}
