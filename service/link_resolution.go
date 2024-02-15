package service

import (
	"context"
	"fmt"
	"github.com/henges/trackrouter/di"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/sync/errgroup"
)

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

	err := eg.Wait()
	return &links, err
}

func resolveLink(text string) {

}
