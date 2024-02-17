package service

import (
	"context"
	"errors"
	"fmt"
	tidalcatalog "github.com/henges/trackrouter/clients/tidal/generate/catalog"
	"github.com/henges/trackrouter/di"
	"github.com/henges/trackrouter/model"
	"github.com/henges/trackrouter/service/helpers"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/sync/errgroup"
	"strings"
)

type LinkResolutionService struct {
	Clients *di.Clients
}

func NewLinkResolutionService(c *di.Clients) *LinkResolutionService {

	return &LinkResolutionService{Clients: c}
}

func (l *LinkResolutionService) FindLinks(message string) (*model.Links, error) {

	id, err := helpers.ResolveId(message)
	if err != nil {
		return nil, err
	}
	metadata, err := l.GetTrackMetadata(id)
	if err != nil {
		return nil, err
	}
	return l.GetLinksFromMetadata(metadata)
}

func (l *LinkResolutionService) GetLinksFromMetadata(md *model.TrackMetadata) (*model.Links, error) {

	q := strings.Join(append([]string{md.Title}, md.Artists...), " ")
	return l.GetLinks(q)
}

func (l *LinkResolutionService) GetLinks(query string) (*model.Links, error) {

	var eg errgroup.Group
	var links model.Links
	ctx := context.TODO()

	eg.Go(func() error {
		search, err := l.Clients.SpotifyClient.Search(ctx, query, spotify.SearchTypeTrack)
		if err != nil {
			return err
		}
		if len(search.Tracks.Tracks) > 0 {
			links.SpotifyLink = fmt.Sprintf("https://open.spotify.com/track/%s", search.Tracks.Tracks[0].ID)
		}
		return nil
	})
	eg.Go(func() error {
		search, err := l.Clients.TidalClient.Search(ctx, query)
		if err != nil {
			return err
		}
		if search.Tracks == nil {
			log.Error().Msg("tracks was nil in Tidal response")
			return nil
		}
		tracks := *search.Tracks
		if len(tracks) > 0 {
			links.TidalLink = tracks[0].Resource.TidalUrl
		}
		return nil
	})
	eg.Go(func() error {
		res, err := l.Clients.YoutubeClient.Search.List([]string{"snippet"}).Q(query).Do()
		if err != nil {
			return err
		}
		if len(res.Items) > 0 {
			links.YoutubeLink = fmt.Sprintf("https://youtube.com/watch?v=%s", res.Items[0].Id.VideoId)
		}
		return nil
	})

	err := eg.Wait()
	return &links, err
}

var ErrUnsupportedProviderType = errors.New("unsupported provider type")

func (l *LinkResolutionService) GetTrackMetadata(id model.ExternalTrackId) (*model.TrackMetadata, error) {

	switch id.ProviderType {
	case model.ProviderTypeSpotify:
		track, err := l.Clients.SpotifyClient.GetTrack(context.TODO(), spotify.ID(id.Id))
		if err != nil {
			return nil, err
		}
		return &model.TrackMetadata{
			Title: track.Name,
			Artists: lo.Map(track.Artists, func(item spotify.SimpleArtist, index int) string {
				return item.Name
			}),
			Album: track.Album.Name,
		}, nil
	case model.ProviderTypeTidal:
		track, err := l.Clients.TidalClient.TrackFromId(context.TODO(), id.Id)
		if err != nil {
			return nil, err
		}
		if track.Data == nil {
			return nil, errors.New("track.Data was nil in tidal response")
		}
		data := *track.Data
		if len(data) <= 0 {
			return nil, errors.New("track.Data was empty in tidal response")
		}
		ft := data[0].Resource
		return &model.TrackMetadata{
			Title: ft.Title,
			Artists: lo.Map(*ft.Artists, func(item tidalcatalog.SimpleArtist, index int) string {
				return item.Name
			}),
			Album: ft.Album.Title,
		}, nil
	}

	return nil, ErrUnsupportedProviderType
}
