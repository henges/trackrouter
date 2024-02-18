package service

import (
	"github.com/henges/trackrouter/model"
	"github.com/henges/trackrouter/providers"
	"github.com/henges/trackrouter/providers/errors"
)

type LinkResolutionService struct {
	Providers providers.Providers
}

func NewLinkResolutionService(ps providers.Providers) *LinkResolutionService {

	return &LinkResolutionService{Providers: ps}
}

func (l *LinkResolutionService) FindLinks(message string) (*model.LinksMatchResult, error) {

	providerType, id, err := l.Providers.MatchId(message)
	if err != nil {
		return nil, err
	}
	metadata, err := l.Providers[providerType].LookupId(id.Id)
	if err != nil {
		return nil, err
	}
	links := l.Providers.Except(providerType).LookupMetadata(metadata)
	if len(links) == 0 {
		return nil, providererrors.ErrIdNotMatched
	}
	return &model.LinksMatchResult{
		Id:    id,
		Links: links,
	}, nil
}
