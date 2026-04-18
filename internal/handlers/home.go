package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handler) HomeHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	fmt.Fprintf(w, "Go server is running and up")
}