package providertesthelpers

import (
	"github.com/henges/trackrouter/model"
	spotifyprovider "github.com/henges/trackrouter/providers/spotify"
	tidalprovider "github.com/henges/trackrouter/providers/tidal"
	"github.com/henges/trackrouter/providers/types"
	youtubeprovider "github.com/henges/trackrouter/providers/youtube"
)

func TestProviders() []providertypes.ProviderMakerFunc {
	return []providertypes.ProviderMakerFunc{
		func() (model.ProviderType, providertypes.Provider) {
			return model.ProviderTypeSpotify, &spotifyprovider.SpotifyProvider{
				SpotifyMatch:  &spotifyprovider.SpotifyMatch{},
				SpotifyLookup: nil,
			}
		}, func() (model.ProviderType, providertypes.Provider) {
			return model.ProviderTypeTidal, &tidalprovider.TidalProvider{
				TidalMatch:  &tidalprovider.TidalMatch{},
				TidalLookup: nil,
			}
		}, func() (model.ProviderType, providertypes.Provider) {
			return model.ProviderTypeYoutube, youtubeprovider.YoutubeProvider{
				YoutubeMatch:  &youtubeprovider.YoutubeMatch{},
				YoutubeLookup: nil,
			}
		},
	}
}
