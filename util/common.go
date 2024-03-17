package util

import (
	"github.com/rs/zerolog/log"
	"sync"
)

func Must[T any](f func() (T, error)) T {

	t, err := f()
	if err != nil {
		log.Fatal().Err(err).Msg("required function returned err")
	}
	return t
}

func ParallelMapValues[K comparable, V1 any, V2 any, M ~map[K]V1](input M, f func(V1) (V2, error)) map[K]V2 {
	var wg sync.WaitGroup
	results := make(map[K]V2, len(input))
	var mu sync.Mutex

	for key, p := range input {
		wg.Add(1)
		key := key
		p := p
		go func() {
			defer wg.Done()
			value, err := f(p)
			if err != nil {
				return
			}
			mu.Lock()
			defer mu.Unlock()
			results[key] = value
		}()
	}

	wg.Wait()
	return results
}
