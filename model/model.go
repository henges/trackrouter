package model

import "fmt"

type LinksMatchResult struct {
	Id    ExternalTrackId
	Links Links
}

type Links map[ProviderType]string

type ProviderType int

const (
	ProviderTypeSpotify ProviderType = 1
	ProviderTypeTidal   ProviderType = 2
	ProviderTypeYoutube ProviderType = 3
)

func (p ProviderType) String() string {
	switch p {
	case ProviderTypeSpotify:
		return "Spotify"
	case ProviderTypeTidal:
		return "Tidal"
	case ProviderTypeYoutube:
		return "Youtube"
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
