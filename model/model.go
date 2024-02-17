package model

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

type TrackMetadata struct {
	Title   string
	Artists []string
	Album   string
}

type Links struct {
	SpotifyLink string `json:"spotifyLink"`
	TidalLink   string `json:"tidalLink"`
	YoutubeLink string `json:"youtubeLink"`
}

func (l Links) IsEmpty() bool {
	return l.SpotifyLink == "" && l.TidalLink == "" && l.YoutubeLink == ""
}
