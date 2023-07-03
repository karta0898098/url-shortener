package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"url-shortener/pkg/app/urlshortener/endpoints"
	"url-shortener/pkg/errors"
)

// Handler is wrap all endpoints
type Handler struct {
	e endpoints.Endpoints
}

// NewHandler new handler
func NewHandler(e endpoints.Endpoints) *Handler {
	return &Handler{e: e}
}

// ShortURL short URL http handler
func (h *Handler) ShortURL(c echo.Context) error {
	var (
		req = new(endpoints.ShortURLRequest)
	)

	if err := c.Bind(req); err != nil {
		return errors.Wrapf(errors.ErrInvalidInput, "failed to bind short url request %v", err)
	}

	if err := c.Validate(req); err != nil {
		return errors.Wrap(errors.ErrInvalidInput, "validate short url request is fail")
	}

	ctx := c.Request().Context()

	resp, err := h.e.ShortURLEndpoint(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

// RedirectURL is redirectURL http handler
func (h *Handler) RedirectURL(c echo.Context) error {
	var (
		req = new(endpoints.RedirectURLRequest)
	)

	if err := c.Bind(req); err != nil {
		return errors.Wrap(errors.ErrInvalidInput, "failed to bind redirect request")
	}

	if err := c.Validate(req); err != nil {
		return errors.Wrap(errors.ErrInvalidInput, "validate redirect url request is fail")
	}

	ctx := c.Request().Context()

	resp, err := h.e.RedirectURLEndpoint(ctx, req)
	if err != nil {
		return err
	}

	r := resp.(*endpoints.RedirectURLResponse)

	return c.Redirect(http.StatusMovedPermanently, r.ShortenedURL.OriginalURL)
}
