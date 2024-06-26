package service

import (
	"github.com/henges/trackrouter/model"
	"github.com/henges/trackrouter/providers"
	"github.com/rs/zerolog/log"
)

type LinkResolutionService struct {
	Providers providers.Providers
	Searchers providers.Searchers
}

func NewLinkResolutionService(ps providers.Providers) *LinkResolutionService {

	return &LinkResolutionService{Providers: ps, Searchers: ps.Searchers()}
}

func (l *LinkResolutionService) FindLinksFromUrl(message string) (*model.UrlMatchResult, error) {

	providerType, id, err := l.Providers.MatchId(message)
	if err != nil {
		return nil, err
	}
	log.Debug().Stringer("providerType", providerType).Str("id", id.Id).Msg("Matched ID")
	metadata, err := l.Providers[providerType].LookupId(id.Id)
	if err != nil {
		return nil, err
	}
	log.Debug().Stringer("providerType", providerType).Str("id", id.Id).Any("metadata", metadata).Msg("Got metadata")
	links := l.Providers.Except(providerType).LookupMetadata(metadata)
	if len(links) == 0 {
		return nil, providers.ErrIdNotMatched
	}
	log.Debug().Stringer("providerType", providerType).
		Str("id", id.Id).Any("metadata", metadata).Any("links", links).Msg("Got metadata")

	return &model.UrlMatchResult{
		Id:            id,
		TrackMetadata: metadata,
		Links:         links,
	}, nil
}

func (l *LinkResolutionService) FindLinksFromMessage(message string) (model.Links, error) {

	links := l.Searchers.LookupQuery(message)
	if len(links) == 0 {
		return nil, providers.ErrIdNotMatched
	}
	log.Debug().Any("links", links).Msg("Got metadata")

	return links, nil
}

func (l *LinkResolutionService) CloneWithJust(providerType model.ProviderType) *LinkResolutionService {

	for k, v := range l.Providers {
		if k == providerType {
			return NewLinkResolutionService(providers.Providers{k: v})
		}
	}
	// Shouldn't be possible I think?
	panic("invalid provider type provided")
}
