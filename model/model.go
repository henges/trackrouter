package model

import "fmt"

type UrlMatchResult struct {
	Id            ExternalTrackId
	TrackMetadata TrackMetadata
	Links         Links
}

type Links map[ProviderType]string

type ProviderType int

const (
	ProviderTypeSpotify ProviderType = 1
	ProviderTypeTidal   ProviderType = 2
	ProviderTypeYoutube ProviderType = 3
	ProviderTypePlex    ProviderType = 4
)

func (p ProviderType) String() string {
	switch p {
	case ProviderTypeSpotify:
		return "Spotify"
	case ProviderTypeTidal:
		return "Tidal"
	case ProviderTypeYoutube:
		return "Youtube"
	case ProviderTypePlex:
		return "Plex"
	}

	return fmt.Sprintf("unknown provider type %d", p)
}

type ExternalTrackId struct {
	ProviderType ProviderType
	Id           string
}

type TrackMetadata struct {
	Title   string
	Artists []string
	Album   string
}
