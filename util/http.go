package util

import (
	"github.com/rs/zerolog/log"
	"net/http"
)

type instrumentedTransport struct {
	delegate http.RoundTripper
}

func (it *instrumentedTransport) RoundTrip(r *http.Request) (*http.Response, error) {

	log.Info().Any("request", r).Send()

	return it.delegate.RoundTrip(r)
}

type InstrumentedTransportOptions struct {
	LogRequests bool
}

func NewInstrumentedTransport(logRequests bool) http.RoundTripper {

	t := http.DefaultTransport.(*http.Transport).Clone()
	if !logRequests {
		return t
	}
	return &instrumentedTransport{delegate: t}
}
