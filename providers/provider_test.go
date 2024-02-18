package providers

import (
	"errors"
	"github.com/henges/trackrouter/providers/errors"
	providertesthelpers "github.com/henges/trackrouter/providers/helpers/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProviders_MatchId(t *testing.T) {

	ps := NewProviders(providertesthelpers.TestProviders()...)
	_, _, err := ps.MatchId("abcd not a link")
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, providererrors.ErrMessageNotMatched))
}
