package helpers

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Success bool      `json:"success"`
	Data    ErrorData `json:"data"`
}

type ErrorData struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func SendJSONError(
	w http.ResponseWriter,
	message string,
	code int,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{
		Success: false,
		Data: ErrorData{
			Error:   http.StatusText(code),
			Message: message,
			Code:    code,
		},
	})
}
