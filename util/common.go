package util

import "github.com/rs/zerolog/log"

func Must[T any](f func() (T, error)) T {

	t, err := f()
	if err != nil {
		log.Fatal().Err(err).Msg("required function returned err")
	}
	return t
}
