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
	Success bool       `json:"success"`
	Data    StreamData `json:"data"`
}

type StreamData struct {
	VideoID string        `json:"video_id"`
	Status  job.JobStatus `json:"status"`
	Message string        `json:"message"`
	Code    int           `json:"code"`
}

func (h *Handler) StreamHandler(w http.ResponseWriter, r *http.Request) {
	
	// 1. Get the video id and actual filename requested (index.m3u8 OR segment.ts)
	videoID := r.PathValue("video_id")
	filename := r.PathValue("filename")

	// 2. handle empty filename and Validate file extension
	if filename == "" {
		filename = "index.m3u8" // Default if only /video_id/ is hit
	} else {
		ext := filepath.Ext(filename)
		if ext != ".m3u8" && ext != ".ts" {
			helpers.SendJSONError(w, "Invalid file type requested", http.StatusBadRequest)
			return
		}
	}

	// 2. Security Check: Prevent Directory Traversal
	// filepath.Clean removes any ".." or redundant slashes
	cleanFilename := filepath.Base(filename)

	outputDir := h.Env.OutputDir // "output"
	videoPath := filepath.Join(outputDir, videoID) // "output/video_id"
	requestedFilePath := filepath.Join(videoPath, cleanFilename) // "output/video_id/{cleanFilename}" 

	// 3. Check if the video directory exists
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		helpers.SendJSONError(w, "Video record not found", http.StatusNotFound)
		return
	}

	// 4. Check if that specific file exists
	if _, err := os.Stat(requestedFilePath); os.IsNotExist(err) {
		if filename == "index.m3u8" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted) // Return 202 Accepted
			json.NewEncoder(w).Encode(StreamResponse{
				Success: true,
				Data: StreamData{
					VideoID: videoID,
					Status:  job.JobStatusProcessing,
					Message: "Master playlist is being generated. Please retry",
					Code:    http.StatusAccepted,
				},
			})
			return
		}

		helpers.SendJSONError(w, "Segment not available yet", http.StatusNotFound)
		return
	}

	// 5. Set Dynamic Content-Type
	if filepath.Ext(cleanFilename) == ".m3u8" {
		w.Header().Set("Content-Type", "application/x-mpegURL")
	} else if filepath.Ext(cleanFilename) == ".ts" {
		w.Header().Set("Content-Type", "video/MP2T")
	}

	// 6. CORS & Serve
	// Note: Since we use CORS middleware in main.go, this might be redundant
	// but keeping it here for safety is fine.
	w.Header().Set("Access-Control-Allow-Origin", "*") // Fixes the CORS loading issue
	w.WriteHeader(http.StatusOK)
	// 7. Serve the specific file (index.m3u8 OR any .ts file)
	http.ServeFile(w, r, requestedFilePath)
}
