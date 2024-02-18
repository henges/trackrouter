package providers_test

import (
	"errors"
	"github.com/henges/trackrouter/di"
	"github.com/henges/trackrouter/providers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProviders_MatchId(t *testing.T) {

	ps := providers.NewProviders(di.TestProviders()...)
	_, _, err := ps.MatchId("abcd not a link")
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, providers.ErrMessageNotMatched))
}
