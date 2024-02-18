package tidalprovider

import (
	"github.com/henges/trackrouter/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatchId_Tidal(t *testing.T) {

	tidalInput := "https://tidal.com/track/82528930 forgot how awesome this is"
	matcher := &TidalMatch{}
	id, err := matcher.MatchId(tidalInput)
	assert.Nil(t, err)
	assert.Equal(t, model.ProviderTypeTidal, id.ProviderType)
	assert.Equal(t, "82528930", id.Id)
}

func TestResolveLink_Tidal_Browse(t *testing.T) {

	tidalInput := "https://tidal.com/browse/track/82528930"
	matcher := &TidalMatch{}
	id, err := matcher.MatchId(tidalInput)
	assert.Nil(t, err)
	assert.Equal(t, model.ProviderTypeTidal, id.ProviderType)
	assert.Equal(t, "82528930", id.Id)
}
