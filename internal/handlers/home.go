package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handler) HomeHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Go server is running and up")
}