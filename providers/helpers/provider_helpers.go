package providerhelpers

import (
	"fmt"
	"github.com/henges/trackrouter/model"
	"github.com/henges/trackrouter/providers/errors"
	"strings"
)

func DefaultNoMatchResult(text string) (model.ExternalTrackId, error) {
	return model.ExternalTrackId{}, fmt.Errorf("for input text %s: %w", text, providererrors.ErrMessageNotMatched)
}

func DefaultTrackMetadataQuery(md model.TrackMetadata) string {
	return strings.Join(append([]string{md.Title}, md.Artists...), " ")
}
