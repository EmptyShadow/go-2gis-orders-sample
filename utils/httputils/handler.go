package httputils

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
)

type Request struct {
	Method  string
	Path    string
	Headers http.Header
	Query   url.Values
	Body    io.Reader
}

type Response struct {
	StatusCode int
	Headers    http.Header
	Body       interface{}
}

var _ http.Handler = (Handler)(nil)

type Handler func(ctx context.Context, req Request) (Response, error)

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp, err := h(r.Context(), Request{
		Method:  r.Method,
		Path:    r.URL.Path,
		Headers: r.Header.Clone(),
		Query:   r.URL.Query(),
		Body:    r.Body,
	})

	var apiErr *APIError
	if !errors.As(err, &apiErr) && err != nil {
		apiErr = NewInternalAPIError()
	}
	if apiErr != nil {
		_ = WriteAPIErrorJSONToResponse(apiErr, w) // TODO: надо бы логировать.
		return
	}

	if resp.StatusCode > 0 {
		w.WriteHeader(resp.StatusCode)
	}
	for key, values := range resp.Headers {
		for _, value := range values {
			w.Header().Set(key, value)
		}
	}
	_ = WriteJSONToResponse(resp.Body, w) // TODO: надо бы логировать.
}

func CheckMethod(method string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		h.ServeHTTP(w, r)
	})
}
