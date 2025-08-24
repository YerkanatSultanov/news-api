package utils

import (
	"encoding/json"
	"net/http"
	"news-api/internal/dto/errors"
)

func WriteError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := errors.ErrorResponse{
		Code:    code,
		Message: message,
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func WriteJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}
