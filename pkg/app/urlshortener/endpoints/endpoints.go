package endpoints

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/spf13/viper"

	"url-shortener/pkg/app/urlshortener/entity"
	"url-shortener/pkg/app/urlshortener/service"
	"url-shortener/pkg/errors"
)

// Endpoints contain all url shortener endpoint
type Endpoints struct {
	ShortURLEndpoint    endpoint.Endpoint
	RedirectURLEndpoint endpoint.Endpoint
}

// New endpoints
func New(svc service.ShortenedURLService) (ep Endpoints) {
	shortURLEndpoint := MakeShortURLEndpoint(svc)
	shortURLEndpoint = endpoint.Chain(
		LoggingMiddleware("shortURL"),
	)(shortURLEndpoint)
	ep.ShortURLEndpoint = shortURLEndpoint

	redirectURLEndpoint := MakeRedirectURLEndpoint(svc)
	redirectURLEndpoint = endpoint.Chain(
		LoggingMiddleware("redirectURL"),
	)(redirectURLEndpoint)
	ep.RedirectURLEndpoint = redirectURLEndpoint

	return ep
}

// ShortURLRequest is define short url request
type ShortURLRequest struct {
	URL       string  `json:"url" validate:"http_url,required"`
	ExpiredAt *string `json:"expireAt"`
}

// ShortURLResponse is define short url response
type ShortURLResponse struct {
	ID       string `json:"id"`
	ShortURL string `json:"shortUrl"`
}

// MakeShortURLEndpoint make short url endpoint
func MakeShortURLEndpoint(svc service.ShortenedURLService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var (
			expireAt *time.Time
		)
		req := request.(*ShortURLRequest)

		if req.ExpiredAt != nil {
			t, err := time.Parse(`2006-01-02T15:04:05Z`, *req.ExpiredAt)
			if err != nil {
				return nil, errors.Wrap(errors.ErrInvalidInput, "input time format is not correct")
			}

			if t.Before(time.Now()) {
				return nil, errors.Wrap(errors.ErrInvalidInput, "input time before now")
			}

			expireAt = &t
		}

		shortenedURL, err := svc.ShortURL(ctx, req.URL, &service.ShortURLOption{
			ExpiredAt: expireAt,
		})
		if err != nil {
			return nil, err
		}

		return &ShortURLResponse{
			ID:       shortenedURL.Short,
			ShortURL: viper.GetString("serverHost") + "/" + shortenedURL.Short,
		}, nil
	}
}

// RedirectURLRequest is redirect short url response
type RedirectURLRequest struct {
	ShortURL string `param:"url"`
}

// RedirectURLResponse is redirect short url response
type RedirectURLResponse struct {
	ShortenedURL *entity.ShortenedURL
}

// MakeRedirectURLEndpoint make redirect url endpoint
func MakeRedirectURLEndpoint(svc service.ShortenedURLService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*RedirectURLRequest)

		shortenedURL, err := svc.RetrieveShortenedURL(ctx, req.ShortURL)
		if err != nil {
			return nil, err
		}

		return &RedirectURLResponse{
			ShortenedURL: shortenedURL,
		}, nil
	}
}
