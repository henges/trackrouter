package tidalprovider

import (
	"context"
	"errors"
	"github.com/henges/trackrouter/clients/tidal"
	tidalcatalog "github.com/henges/trackrouter/clients/tidal/generate/catalog"
	"github.com/henges/trackrouter/model"
	"github.com/henges/trackrouter/providers/helpers"
	"github.com/henges/trackrouter/providers/types"
	"github.com/henges/trackrouter/util"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"regexp"
)

func NewTidalProvider(c tidal.Client) providertypes.Provider {

	return &TidalProvider{
		TidalMatch:  &TidalMatch{},
		TidalLookup: &TidalLookup{c},
	}
}

type TidalProvider struct {
	*TidalMatch
	*TidalLookup
}

type TidalMatch struct{}

type TidalLookup struct {
	client tidal.Client
}

var tidalRegex = regexp.MustCompile("tidal\\.com/(browse/)?track/(\\d+)")

func (s *TidalMatch) MatchId(text string) (model.ExternalTrackId, error) {

	if match := util.RegexpMatchWithGroup(text, tidalRegex); match != "" {
		return model.ExternalTrackId{ProviderType: model.ProviderTypeTidal, Id: match}, nil
	}

	return providerhelpers.DefaultNoMatchResult(text)
}

func (s *TidalLookup) LookupId(id string) (model.TrackMetadata, error) {

	track, err := s.client.TrackFromId(context.TODO(), id)
	if err != nil {
		return model.TrackMetadata{}, err
	}
	if track.Data == nil {
		return model.TrackMetadata{}, errors.New("track.Data was nil in tidal response")
	}
	data := *track.Data
	if len(data) <= 0 {
		return model.TrackMetadata{}, errors.New("track.Data was empty in tidal response")
	}
	ft := data[0].Resource
	if ft.Artists == nil {
		return model.TrackMetadata{}, errors.New("track.Data.Artists was nil in tidal response")
	}
	return model.TrackMetadata{
		Title: ft.Title,
		Artists: lo.Map(*ft.Artists, func(item tidalcatalog.SimpleArtist, index int) string {
			return item.Name
		}),
		Album: ft.Album.Title,
	}, nil
}

func (s *TidalLookup) LookupMetadata(metadata model.TrackMetadata) string {

	query := providerhelpers.DefaultTrackMetadataQuery(metadata)
	search, err := s.client.Search(context.Background(), query)
	if err != nil {
		log.Error().Err(err).Msg("in tidal request")
		return ""
	}
	if search.Tracks == nil {
		log.Error().Msg("tracks was nil in Tidal response")
		return ""
	}
	tracks := *search.Tracks
	if len(tracks) > 0 {
		return tracks[0].Resource.TidalUrl
	}
	return ""
}
