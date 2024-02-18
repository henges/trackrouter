package di

import (
	"github.com/henges/trackrouter/model"
	"github.com/henges/trackrouter/providers"
	"github.com/henges/trackrouter/providers/plex"
	"github.com/henges/trackrouter/providers/spotify"
	"github.com/henges/trackrouter/providers/tidal"
	"github.com/henges/trackrouter/providers/youtube"
)

func TestProviders() []providers.ProviderMakerFunc {
	return []providers.ProviderMakerFunc{
		func() (model.ProviderType, providers.Provider) {
			return model.ProviderTypeSpotify, &spotifyprovider.SpotifyProvider{
				SpotifyMatch:  &spotifyprovider.SpotifyMatch{},
				SpotifyLookup: nil,
			}
		}, func() (model.ProviderType, providers.Provider) {
			return model.ProviderTypeTidal, &tidalprovider.TidalProvider{
				TidalMatch:  &tidalprovider.TidalMatch{},
				TidalLookup: nil,
			}
		}, func() (model.ProviderType, providers.Provider) {
			return model.ProviderTypeYoutube, youtubeprovider.YoutubeProvider{
				YoutubeMatch:  &youtubeprovider.YoutubeMatch{},
				YoutubeLookup: nil,
			}
		}, func() (model.ProviderType, providers.Provider) {
			return model.ProviderTypePlex, plexprovider.NewPlexProvider()
		},
	}
}
