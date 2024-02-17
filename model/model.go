package model

import "fmt"

type LinksMatchResult struct {
	Id    ExternalTrackId
	Links *Links
}

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

type Links struct {
	SpotifyLink string `json:"spotifyLink"`
	TidalLink   string `json:"tidalLink"`
	YoutubeLink string `json:"youtubeLink"`
}

func (l Links) Count() int {
	count := 0
	if l.SpotifyLink != "" {
		count++
	}
	if l.TidalLink != "" {
		count++
	}
	if l.YoutubeLink != "" {
		count++
	}
	return count
}

func (l Links) IsEmpty() bool {
	return l.SpotifyLink == "" && l.TidalLink == "" && l.YoutubeLink == ""
}
