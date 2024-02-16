package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"sync"
)

type Config struct {
	Spotify *SpotifyConfig
	Logging *LoggingConfig
	Tidal   *TidalConfig
}

type LoggingConfig struct {
	Level string `split_words:"true" default:"trace"`
	Style string `split_words:"true" default:"friendly"`
}

type SpotifyConfig struct {
	LogRequests  bool   `split_words:"true"`
	ClientId     string `split_words:"true"`
	ClientSecret string `split_words:"true"`
}

type TidalConfig struct {
	LogRequests  bool   `split_words:"true"`
	ClientId     string `split_words:"true"`
	ClientSecret string `split_words:"true"`
}

var c Config
var once sync.Once

func Get() *Config {

	once.Do(func() {
		envconfig.MustProcess("", &c)
		configureLogging(c.Logging)
		log.Info().Any("config", c).Send()
	})
	return &c
}

func configureLogging(c *LoggingConfig) {

	level, err := zerolog.ParseLevel(c.Level)
	if err != nil {
		log.Fatal().Err(err).Msg("error marshalling log level")
	}
	log.Logger = log.Logger.Level(level)

	switch c.Style {
	case "friendly":
		log.Logger = log.Logger.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		break
		// use the json default
	}
}
