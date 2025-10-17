package handlers

import (
	"encoding/json"
	"net/http"
)

// Resp is the standardized API response envelope.
type Resp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// WriteSuccess writes a successful response with code=0 and message="success".
func WriteSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(Resp{Code: 0, Message: "success", Data: data})
}

// WriteError writes an error response with provided code and message. HTTP status
// is always 200 per API contract; code != 0 signals failure.
func WriteError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(Resp{Code: code, Message: message, Data: nil})
}
