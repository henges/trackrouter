package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolveLink_Spotify(t *testing.T) {

	spotifyInput := "https://open.spotify.com/track/64Yjy4alpDOtrJuzIcS8O5?si=nvLdrfLNTEO62ZKxCxdJDw"
	id, err := ResolveId(spotifyInput)
	assert.Nil(t, err)
	assert.Equal(t, ProviderTypeSpotify, id.ProviderType)
	assert.Equal(t, "64Yjy4alpDOtrJuzIcS8O5", id.Id)
}

func TestResolveLink_Tidal(t *testing.T) {

	tidalInput := "https://tidal.com/track/82528930 forgot how awesome this is"
	id, err := ResolveId(tidalInput)
	assert.Nil(t, err)
	assert.Equal(t, ProviderTypeTidal, id.ProviderType)
	assert.Equal(t, "82528930", id.Id)
}
