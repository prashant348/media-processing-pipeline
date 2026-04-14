package handlers

import "net/http"

func VideoHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "video.html")
}