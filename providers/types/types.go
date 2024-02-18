package providertypes

import "github.com/henges/trackrouter/model"

type Provider interface {
	// MatchId parses text and tries to locate an ID for a track
	// that originates from this provider. It should return a wrapped
	// ErrMessageNotMatched when no match is found.
	MatchId(text string) (model.ExternalTrackId, error)

	// LookupId attempts to locate track metadata based on an ID.
	LookupId(id string) (model.TrackMetadata, error)

	// LookupMetadata attempts to locate a matching ID for a track
	// based on the input metadata.
	LookupMetadata(metadata model.TrackMetadata) string
}

type ProviderMakerFunc func() (model.ProviderType, Provider)
