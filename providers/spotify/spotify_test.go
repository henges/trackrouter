package spotifyprovider

import (
	"github.com/henges/trackrouter/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatchId_Spotify(t *testing.T) {

	spotifyInput := "https://open.spotify.com/track/64Yjy4alpDOtrJuzIcS8O5?si=nvLdrfLNTEO62ZKxCxdJDw"
	matcher := SpotifyMatch{}
	id, err := matcher.MatchId(spotifyInput)
	assert.Nil(t, err)
	assert.Equal(t, model.ProviderTypeSpotify, id.ProviderType)
	assert.Equal(t, "64Yjy4alpDOtrJuzIcS8O5", id.Id)
}
