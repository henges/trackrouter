package service

import (
	"context"
	"errors"
	"fmt"
	tidalcatalog "github.com/henges/trackrouter/clients/tidal/generate/catalog"
	"github.com/henges/trackrouter/di"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/sync/errgroup"
	"regexp"
	"strings"
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

func (l *LinkResolutionService) GetLinksFromMetadata(md *TrackMetadata) (*Links, error) {

	q := strings.Join(append([]string{md.Title}, md.Artists...), " ")
	return l.GetLinks(q)
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
			log.Error().Msg("tracks was nil in Tidal response")
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

type ProviderType int

const (
	ProviderTypeSpotify ProviderType = 1
	ProviderTypeTidal   ProviderType = 2
	ProviderTypeYoutube ProviderType = 3
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

type TrackMetadata struct {
	Title   string
	Artists []string
	Album   string
}

var ErrUnsupportedProviderType = errors.New("unsupported provider type")

func (l *LinkResolutionService) GetTrackMetadata(id ExternalTrackId) (*TrackMetadata, error) {

	switch id.ProviderType {
	case ProviderTypeSpotify:
		track, err := l.Clients.SpotifyClient.GetTrack(context.TODO(), spotify.ID(id.Id))
		if err != nil {
			return nil, err
		}
		return &TrackMetadata{
			Title: track.Name,
			Artists: lo.Map(track.Artists, func(item spotify.SimpleArtist, index int) string {
				return item.Name
			}),
			Album: track.Album.Name,
		}, nil
	case ProviderTypeTidal:
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
		return &TrackMetadata{
			Title: ft.Title,
			Artists: lo.Map(*ft.Artists, func(item tidalcatalog.SimpleArtist, index int) string {
				return item.Name
			}),
			Album: ft.Album.Title,
		}, nil
	}

	return nil, ErrUnsupportedProviderType
}
