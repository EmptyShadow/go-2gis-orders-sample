package httputils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteJSONToResponse(bodyObj interface{}, w http.ResponseWriter) error {
	w.Header().Add(ContentType, ApplicationJSON)

	body, err := json.Marshal(bodyObj)
	if err != nil {
		return fmt.Errorf("json marshal body: %w", err)
	}
	if _, err = w.Write(body); err != nil {
		return fmt.Errorf("write body to response: %w", err)
	}

	return nil
}
