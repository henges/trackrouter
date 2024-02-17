package helpers

import (
	"errors"
	"github.com/henges/trackrouter/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolveLink_Spotify(t *testing.T) {

	spotifyInput := "https://open.spotify.com/track/64Yjy4alpDOtrJuzIcS8O5?si=nvLdrfLNTEO62ZKxCxdJDw"
	id, err := ResolveId(spotifyInput)
	assert.Nil(t, err)
	assert.Equal(t, model.ProviderTypeSpotify, id.ProviderType)
	assert.Equal(t, "64Yjy4alpDOtrJuzIcS8O5", id.Id)
}

func TestResolveLink_Tidal(t *testing.T) {

	tidalInput := "https://tidal.com/track/82528930 forgot how awesome this is"
	id, err := ResolveId(tidalInput)
	assert.Nil(t, err)
	assert.Equal(t, model.ProviderTypeTidal, id.ProviderType)
	assert.Equal(t, "82528930", id.Id)
}

func TestResolveLink_NoMatch_Error(t *testing.T) {

	input := "abcd not a link"
	_, err := ResolveId(input)
	assert.True(t, errors.Is(err, ErrNoMatch))
}
