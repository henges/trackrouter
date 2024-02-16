package tidal

import (
	"context"
	"fmt"
	tidalgen "github.com/henges/trackrouter/clients/tidal/generate"
	"github.com/henges/trackrouter/config"
	"github.com/henges/trackrouter/util"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"strings"
)

type Client interface {
	Search(ctx context.Context, query string) (*SearchResponse, error)
}

type SearchParams = tidalgen.SearchParams
type SearchResponse = tidalgen.SearchResponse

type Error struct {
	statusCode int
	errors     []tidalgen.Error
}

func (e *Error) Error() string {
	if e.errors == nil {
		return fmt.Sprintf("status: %v; errors was nil", e.statusCode)
	}

	var sb strings.Builder
	for _, err := range e.errors {
		var detail, field string
		if err.Detail != nil {
			detail = *err.Detail
		}
		if err.Field != nil {
			field = *err.Field
		}
		sb.WriteString(fmt.Sprintf("status: %v, category: %v, code: %v, detail: %v, field: %v", e.statusCode, err.Category, err.Code, detail, field))
	}

	return sb.String()
}

type client struct {
	delegate tidalgen.ClientWithResponsesInterface
}

func NewClient(c *config.TidalConfig) (Client, error) {
	auth := clientcredentials.Config{
		ClientID:     c.ClientId,
		ClientSecret: c.ClientSecret,
		TokenURL:     "https://auth.tidal.com/v1/oauth2/token",
	}
	ctx := context.WithValue(context.TODO(),
		oauth2.HTTPClient,
		&http.Client{Transport: util.NewInstrumentedTransport(c.LogRequests)})
	httpClient := auth.Client(ctx)

	delegate, err := tidalgen.NewClientWithResponses("https://openapi.tidal.com", tidalgen.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}
	return &client{delegate}, nil
}

func (c *client) Search(ctx context.Context, query string) (*SearchResponse, error) {

	offset, limit, countryCode, popularity := "0", "10", "US", "WORLDWIDE"
	params := SearchParams{
		Query:       query,
		Offset:      &offset,
		Limit:       &limit,
		CountryCode: countryCode,
		Popularity:  (*tidalgen.SearchParamsPopularity)(&popularity),
	}
	response, err := c.delegate.SearchWithBodyWithResponse(
		ctx,
		&params,
		// Need to send this content type with empty body, otherwise their WAF rejects it
		"application/vnd.tidal.v1+json",
		strings.NewReader(""))
	if err != nil {
		return nil, err
	}

	status := response.StatusCode()
	if status >= 200 && status < 300 {
		return response.ApplicationvndTidalV1JSON207, nil
	}
	var errType *tidalgen.ErrorsResponse
	switch status {
	case 400:
		errType = response.ApplicationvndTidalV1JSON400
		break
	case 403:
		body := string(response.Body)
		return nil, &Error{statusCode: status, errors: []tidalgen.Error{{Detail: &body}}}
	case 404:
		errType = response.ApplicationvndTidalV1JSON404
		break
	case 405:
		errType = response.ApplicationvndTidalV1JSON405
		break
	case 406:
		errType = response.ApplicationvndTidalV1JSON406
		break
	case 415:
		errType = response.ApplicationvndTidalV1JSON415
		break
	case 500:
		errType = response.ApplicationvndTidalV1JSON500
		break
	}
	if errType == nil {
		return nil, &Error{status, nil}
	}

	return nil, &Error{status, errType.Errors}
}
