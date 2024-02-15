package config

import (
	"github.com/kelseyhightower/envconfig"
	"sync"
)

type Config struct {
	Spotify *SpotifyConfig
}

type SpotifyConfig struct {
	ClientId     string `split_words:"true"`
	ClientSecret string `split_words:"true"`
}

var c *Config
var once sync.Once

func Get() *Config {

	once.Do(func() {
		envconfig.MustProcess("", c)
	})
	return c
}
