package util

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

type instrumentedTransport struct {
	delegate http.RoundTripper
}

func (it *instrumentedTransport) RoundTrip(r *http.Request) (*http.Response, error) {

	marshalled := marshalRequest(r)
	log.Info().Msg(marshalled)

	return it.delegate.RoundTrip(r)
}

func marshalRequest(r *http.Request) string {

	var sb []string
	sb = append(sb, fmt.Sprintf("%s %s %v", r.Proto, r.Method, r.URL))
	for k, v := range r.Header {
		for _, v2 := range v {
			sb = append(sb, fmt.Sprintf("%s: %s", k, v2))
		}
	}

	return strings.Join(sb, "\n")
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
