package handlers

import "net/http"

func (h *Handler) VideoHandler(
	w http.ResponseWriter, 
	r *http.Request,
) {
	filePath := "web/templates/video.html"
	http.ServeFile(w, r, filePath)
}