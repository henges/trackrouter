package di

import (
	"context"
	"github.com/henges/trackrouter/config"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
	"sync"
)

type Clients struct {
	SpotifyClient *spotify.Client
}

type Deps struct {
	Clients *Clients
}

func mustInitialiseSpotify(c *config.SpotifyConfig) *spotify.Client {

	auth := clientcredentials.Config{
		ClientID:     c.ClientId,
		ClientSecret: c.ClientSecret,
		TokenURL:     spotifyauth.TokenURL,
		Scopes:       []string{"clientcredentials"},
	}
	httpClient := auth.Client(context.Background())

	return spotify.New(httpClient)
}

func mustInitialise(c *config.Config) *Deps {

	d := &Deps{Clients: &Clients{}}
	d.Clients.SpotifyClient = mustInitialiseSpotify(c.Spotify)

	return d
}

var once sync.Once
var deps *Deps

func Get(c *config.Config) *Deps {

	once.Do(func() {
		deps = mustInitialise(c)
	})
	return deps
}
