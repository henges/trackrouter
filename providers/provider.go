package providers

import (
	"github.com/henges/trackrouter/model"
	"github.com/henges/trackrouter/providers/helpers"
	"github.com/henges/trackrouter/providers/types"
	"sync"
)

func NewProviders(providerFuncs ...providertypes.ProviderMakerFunc) Providers {
	ret := make(Providers, len(providerFuncs))
	for _, f := range providerFuncs {
		t, p := f()
		ret[t] = p
	}

	return ret
}

func (ps Providers) MatchId(text string) (model.ProviderType, model.ExternalTrackId, error) {

	for t, p := range ps {
		m, err := p.MatchId(text)
		if err == nil {
			return t, m, nil
		}
	}
	none, err := providerhelpers.DefaultNoMatchResult(text)
	return 0, none, err
}

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

type Providers map[model.ProviderType]providertypes.Provider
