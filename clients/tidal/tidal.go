package tidal

import (
	"context"
	"fmt"
	"github.com/henges/trackrouter/clients/tidal/generate/catalog"
	"github.com/henges/trackrouter/clients/tidal/generate/search"
	"github.com/henges/trackrouter/config"
	"github.com/henges/trackrouter/util"
	"github.com/samber/lo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"strings"
)

type Client interface {
	Search(ctx context.Context, query string) (*SearchResponse, error)
	TrackFromId(ctx context.Context, trackId string) (*tidalcatalog.TracksMultiStatusResponse, error)
}

type SearchParams = tidalsearch.SearchParams
type SearchResponse = tidalsearch.SearchResponse

type tidalError struct {
	statusCode int
	errors     []tidalErrorDetail
}

type tidalErrorDetail struct {
	Category string
	Code     string
	Detail   string
	Field    string
}

func newTidalSearchError(statusCode int, errs []tidalsearch.Error) *tidalError {

	errDetails := lo.Map(errs, func(item tidalsearch.Error, index int) tidalErrorDetail {

		var detail, field string
		if item.Detail != nil {
			detail = *item.Detail
		}
		if item.Field != nil {
			field = *item.Field
		}
		return tidalErrorDetail{
			Category: string(item.Category),
			Code:     item.Code,
			Detail:   detail,
			Field:    field,
		}
	})

	return &tidalError{statusCode: statusCode, errors: errDetails}
}

func newTidalCatalogError(statusCode int, errs []tidalcatalog.Error) *tidalError {

	errDetails := lo.Map(errs, func(item tidalcatalog.Error, index int) tidalErrorDetail {

		var detail, field string
		if item.Detail != nil {
			detail = *item.Detail
		}
		if item.Field != nil {
			field = *item.Field
		}
		return tidalErrorDetail{
			Category: string(item.Category),
			Code:     item.Code,
			Detail:   detail,
			Field:    field,
		}
	})

	return &tidalError{statusCode: statusCode, errors: errDetails}
}

func (e *tidalError) Error() string {
	if e.errors == nil {
		return fmt.Sprintf("status: %v; errors was nil", e.statusCode)
	}

	var sb strings.Builder
	for _, err := range e.errors {
		sb.WriteString(fmt.Sprintf("status: %v, category: %v, code: %v, detail: %v, field: %v", e.statusCode, err.Category, err.Code, err.Detail, err.Field))
	}

	return sb.String()
}

type client struct {
	searcher tidalsearch.ClientWithResponsesInterface
	catalog  tidalcatalog.ClientWithResponsesInterface
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

	searcher, err := tidalsearch.NewClientWithResponses("https://openapi.tidal.com", tidalsearch.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}
	catalog, err := tidalcatalog.NewClientWithResponses("https://openapi.tidal.com", tidalcatalog.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}
	return &client{searcher: searcher, catalog: catalog}, nil
}

func (c *client) Search(ctx context.Context, query string) (*SearchResponse, error) {

	offset, limit, countryCode, popularity := "0", "10", "US", "WORLDWIDE"
	params := SearchParams{
		Query:       query,
		Offset:      &offset,
		Limit:       &limit,
		CountryCode: countryCode,
		Popularity:  (*tidalsearch.SearchParamsPopularity)(&popularity),
	}
	response, err := c.searcher.SearchWithBodyWithResponse(
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
	var errType *tidalsearch.ErrorsResponse
	switch status {
	case 400:
		errType = response.ApplicationvndTidalV1JSON400
		break
	case 403:
		return nil, &tidalError{statusCode: status, errors: []tidalErrorDetail{{Detail: string(response.Body)}}}
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
		return nil, &tidalError{status, nil}
	}

	return nil, newTidalSearchError(status, errType.Errors)
}

func (c *client) TrackFromId(ctx context.Context, trackId string) (*tidalcatalog.TracksMultiStatusResponse, error) {

	params := tidalcatalog.GetTracksByIdsParams{
		Ids:         []string{trackId},
		CountryCode: "US",
	}
	response, err := c.catalog.GetTracksByIdsWithBodyWithResponse(ctx,
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
	var errType *tidalcatalog.ErrorsResponse
	switch status {
	case 400:
		errType = response.ApplicationvndTidalV1JSON400
		break
	case 403:
		return nil, &tidalError{statusCode: status, errors: []tidalErrorDetail{{Detail: string(response.Body)}}}
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
		return nil, &tidalError{status, nil}
	}

	return nil, newTidalCatalogError(status, errType.Errors)
}
