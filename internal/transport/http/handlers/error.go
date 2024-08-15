package handlers

import (
	"encoding/json"
	"net/http"
)

type ErrorMessage struct {
	StatusCode int    `json:"status_code,omitempty"`
	Message    string `json:"message,omitempty"`
	URL        string `json:"url,omitempty"`
}

func WriteErrorResponse(w http.ResponseWriter, r *http.Request, err error, statusCode int) {
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(ErrorMessage{
		StatusCode: statusCode,
		Message:    err.Error(),
		URL:        r.URL.Host + r.URL.Path,
	})
}
