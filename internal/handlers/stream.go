package handlers

import (
	"encoding/json"
	"media_processing_pipeline/internal/helpers"
	"media_processing_pipeline/internal/job"
	"net/http"
	"os"
	"path/filepath"
)

type StreamResponse struct {
	VideoID string        `json:"video_id"`
	Status  job.JobStatus `json:"status"`
	Message string        `json:"message"`
}

func (h *Handler) StreamHandler(w http.ResponseWriter, r *http.Request) {

	videoID := r.PathValue("video_id")
	filename := "index.m3u8"

	outputDir := h.Env.OutputDir
	videoPath := filepath.Join(outputDir, videoID)

	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		helpers.SendJSONError(w, "Video not found or has not started processing yet", http.StatusNotFound)
		return 
	}

	playlistPath := filepath.Join(videoPath, filename)
	if _, err := os.Stat(playlistPath); os.IsNotExist(err) {
		json.NewEncoder(w).Encode(StreamResponse{
			VideoID: videoID,
			Status: "processing",
			Message: "Video is being transcoded, please try again in a few seconds",
		})
		return
	}

	http.ServeFile(w, r, playlistPath)
}
