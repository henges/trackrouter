package di

import (
	"context"
	"github.com/henges/trackrouter/clients/tidal"
	"github.com/henges/trackrouter/config"
	"github.com/henges/trackrouter/model"
	"github.com/henges/trackrouter/providers/spotify"
	"github.com/henges/trackrouter/providers/tidal"
	"github.com/henges/trackrouter/providers/types"
	"github.com/henges/trackrouter/providers/youtube"
	"github.com/henges/trackrouter/util"
	"github.com/rs/zerolog/log"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"net/http"
	"sync"
)

type Clients struct {
	SpotifyClient *spotify.Client
	TidalClient   tidal.Client
	YoutubeClient *youtube.Service
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

func mustInitialiseYoutube(c *config.YoutubeConfig) *youtube.Service {

	svc, err := youtube.NewService(context.TODO(),
		//option.WithHTTPClient(&http.Client{Transport: util.NewInstrumentedTransport(c.LogRequests)}),
		option.WithAPIKey(c.ApiKey))
	if err != nil {
		log.Fatal().Err(err).Msg("while initialising Youtube client")
	}
	return svc
}

func mustInitialise(c *config.Config) *Deps {

	d := &Deps{Clients: &Clients{}}
	d.Clients.SpotifyClient = mustInitialiseSpotify(c.Spotify)
	d.Clients.TidalClient = mustInitialiseTidal(c.Tidal)
	d.Clients.YoutubeClient = mustInitialiseYoutube(c.Youtube)

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

func DefaultProvidersFromDeps(clients *Clients) []providertypes.ProviderMakerFunc {

	return []providertypes.ProviderMakerFunc{
		func() (model.ProviderType, providertypes.Provider) {
			return model.ProviderTypeSpotify, spotifyprovider.NewSpotifyProvider(clients.SpotifyClient)
		}, func() (model.ProviderType, providertypes.Provider) {
			return model.ProviderTypeTidal, tidalprovider.NewTidalProvider(clients.TidalClient)
		}, func() (model.ProviderType, providertypes.Provider) {
			return model.ProviderTypeYoutube, youtubeprovider.NewYoutubeProvider(clients.YoutubeClient)
		},
	}
}
