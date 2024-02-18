package youtubeprovider

import (
	"fmt"
	"github.com/henges/trackrouter/model"
	"github.com/henges/trackrouter/providers"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/youtube/v3"
	"regexp"
)

func NewYoutubeProvider(c *youtube.Service) providers.Provider {

	return &YoutubeProvider{
		YoutubeMatch:  &YoutubeMatch{},
		YoutubeLookup: &YoutubeLookup{c},
	}
}

type YoutubeProvider struct {
	*YoutubeMatch
	*YoutubeLookup
}

type YoutubeMatch struct{}

type YoutubeLookup struct {
	client *youtube.Service
}

var youtubeRegex = regexp.MustCompile("youtube\\.com/watch?v=([^&]+)")

func (s *YoutubeMatch) MatchId(text string) (model.ExternalTrackId, error) {
	return model.ExternalTrackId{}, providers.ErrUnsupportedOperations
	//if match := util.RegexpMatchWithGroup(text, youtubeRegex); match != "" {
	//	return model.ExternalTrackId{ProviderType: model.ProviderTypeYoutube, Id: match}, nil
	//}
	//
	//return providers.DefaultNoMatchResult(text)
}

func (s *YoutubeLookup) LookupId(id string) (model.TrackMetadata, error) {

	return model.TrackMetadata{}, providers.ErrUnsupportedOperations
}

func (s *YoutubeLookup) LookupMetadata(metadata model.TrackMetadata) string {

	query := providers.DefaultTrackMetadataQuery(metadata)
	res, err := s.client.Search.List([]string{"snippet"}).Q(query).Do()
	if err != nil {
		log.Error().Err(err).Msg("in youtube request")
		return ""
	}
	if len(res.Items) > 0 {
		return fmt.Sprintf("https://youtube.com/watch?v=%s", res.Items[0].Id.VideoId)
	}
	return ""
}
