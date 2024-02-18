package providers

import (
	"errors"
	"fmt"
	"github.com/henges/trackrouter/model"
	"strings"
	"sync"
)

type Providers map[model.ProviderType]Provider

type Provider interface {
	// MatchId parses text and tries to locate an ID for a track
	// that originates from this provider. It should return a wrapped
	// ErrMessageNotMatched when no match is found.
	MatchId(text string) (model.ExternalTrackId, error)

	// LookupId attempts to locate track metadata based on an ID.
	LookupId(id string) (model.TrackMetadata, error)

	// LookupMetadata attempts to locate a matching URL for a track
	// based on the input metadata.
	LookupMetadata(metadata model.TrackMetadata) string
}

type ProviderMakerFunc func() (model.ProviderType, Provider)

var ErrMessageNotMatched = errors.New("no match found for input text")
var ErrIdNotMatched = errors.New("no other provider had a track that matched")
var ErrUnsupportedOperations = errors.New("unsupported operation")

func NewProviders(providerFuncs ...ProviderMakerFunc) Providers {
	ret := make(Providers, len(providerFuncs))
	for _, f := range providerFuncs {
		t, p := f()
		ret[t] = p
	}

	return ret
}

// MatchId executes [providertypes.Provider.MatchId] for each providertypes.Provider, returning
// a model.ExternalTrackId for the first non-error result and the model.ProviderType of the
// provider that answered with that result.
func (ps Providers) MatchId(text string) (model.ProviderType, model.ExternalTrackId, error) {

	var id model.ExternalTrackId
	var err error
	for t, p := range ps {
		id, err = p.MatchId(text)
		if err == nil {
			return t, id, nil
		}
	}
	return 0, id, err
}

// Except returns a shallow copy of this Providers minus the provider
// with the specified model.ProviderType.
func (ps Providers) Except(providerType model.ProviderType) Providers {

	ret := make(Providers, len(ps)-1)
	for k, v := range ps {
		if k != providerType {
			ret[k] = v
		}
	}

	return ret
}

func (ps Providers) LookupMetadata(metadata model.TrackMetadata) model.Links {

	var wg sync.WaitGroup
	results := make(model.Links, len(ps))
	var mu sync.Mutex

	for providerType, p := range ps {
		wg.Add(1)
		providerType := providerType
		p := p
		go func() {
			defer wg.Done()
			match := p.LookupMetadata(metadata)
			if match == "" {
				return
			}
			mu.Lock()
			defer mu.Unlock()
			results[providerType] = match
		}()
	}

	wg.Wait()
	return results
}

func DefaultNoMatchResult(text string) (model.ExternalTrackId, error) {
	return model.ExternalTrackId{}, fmt.Errorf("for input text %s: %w", text, ErrMessageNotMatched)
}

func DefaultTrackMetadataQuery(md model.TrackMetadata) string {
	return strings.Join(append([]string{md.Title}, md.Artists...), " ")
}
