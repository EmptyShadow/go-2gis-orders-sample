package httputils

import (
	"errors"
	"fmt"
	"net/http"
)

func WriteAPIErrorJSONToResponse(apiErr *APIError, w http.ResponseWriter) error {
	w.WriteHeader(apiErr.StatusCode)

	if err := WriteJSONToResponse(apiErr, w); err != nil {
		return fmt.Errorf("write json to response: %w", err)
	}

	return nil
}

type APIError struct {
	StatusCode int           `json:"-"`
	Message    string        `json:"message"`
	Details    []interface{} `json:"details,omitempty"`
}

func NewAPIError(code int, message string, details ...interface{}) error {
	return &APIError{
		StatusCode: code,
		Message:    message,
		Details:    details,
	}
}

func APIErrorFromError(err error) (*APIError, bool) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr, true
	}
	return nil, false
}

func (e *APIError) Error() string {
	return fmt.Sprintf("StatusCode='%s' Message='%s'", http.StatusText(e.StatusCode), e.Message)
}

func NewInternalAPIError() *APIError {
	return &APIError{
		StatusCode: http.StatusInternalServerError,
		Message:    "internal",
	}
}

func NewBadArgumentsAPIError(details ...BadArgument) *APIError {
	_details := make([]interface{}, len(details))
	for i, detail := range details {
		_details[i] = detail
	}

	return &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    "bad arguments",
		Details:    _details,
	}
}

type BadArgument struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

func NewInvalidBodyFormatAPIError() *APIError {
	return &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    "invalid body format",
	}
}
