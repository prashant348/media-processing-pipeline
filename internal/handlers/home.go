package handlers

import "net/http"

func (h *Handler) HomeHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	filePath := "web/templates/index.html"
	http.ServeFile(w, r, filePath)
}