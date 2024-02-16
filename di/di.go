package di

import (
	"context"
	"github.com/henges/trackrouter/clients/tidal"
	"github.com/henges/trackrouter/config"
	"github.com/henges/trackrouter/util"
	"github.com/rs/zerolog/log"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"sync"
)

type Clients struct {
	SpotifyClient *spotify.Client
	TidalClient   tidal.Client
}

type Deps struct {
	Clients *Clients
}

func mustInitialiseSpotify(c *config.SpotifyConfig) *spotify.Client {

	auth := clientcredentials.Config{
		ClientID:     c.ClientId,
		ClientSecret: c.ClientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}
	ctx := context.WithValue(context.TODO(),
		oauth2.HTTPClient,
		&http.Client{Transport: util.NewInstrumentedTransport(c.LogRequests)})
	httpClient := auth.Client(ctx)

	return spotify.New(httpClient)
}

func mustInitialiseTidal(c *config.TidalConfig) tidal.Client {

	client, err := tidal.NewClient(c)
	if err != nil {
		log.Fatal().Err(err).Msg("while initialising Tidal client")
	}
	return client
}

func mustInitialise(c *config.Config) *Deps {

	d := &Deps{Clients: &Clients{}}
	d.Clients.SpotifyClient = mustInitialiseSpotify(c.Spotify)
	d.Clients.TidalClient = mustInitialiseTidal(c.Tidal)

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
