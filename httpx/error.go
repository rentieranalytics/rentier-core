package httpx

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code    int                 `json:"-"`
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors,omitempty"`
}

func WriteError(w http.ResponseWriter, status int, message string, errs map[string][]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := ErrorResponse{
		Code:    status,
		Message: message,
		Errors:  errs,
	}
	_ = json.NewEncoder(w).Encode(resp)
}
