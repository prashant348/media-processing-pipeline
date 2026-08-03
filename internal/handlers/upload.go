package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"log"
	"media_processing_pipeline/internal/helpers"
	"media_processing_pipeline/internal/job"

	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type UploadResponse struct {
	Success bool       `json:"success"`
	Data    UploadData `json:"data"`
}

type UploadData struct {
	VideoID   string        `json:"video_id"`
	JobID     string        `json:"job_id"`
	Status    job.JobStatus `json:"status"`
	Metadata  FileMetaData  `json:"metadata"`
	Links     ActionLinks   `json:"links"`
	CreatedAt time.Time     `json:"created_at"`
}

type FileMetaData struct {
	Filename    string `json:"filename"`
	SizeBytes   int64  `json:"size_bytes"`
	ContentType string `json:"content_type"`
}

type ActionLinks struct {
	StatusURL string `json:"status_url"`
	StreamURL string `json:"stream_url"`
	VideoURL  string `json:"video_url"`
}

func (h *Handler) UploadHandler(w http.ResponseWriter, r *http.Request) {
	// set a memory limit of 100KB
	const maxMemory = 100 << 10
	// this method is used to process/parse "multipart/form-data" MIME type request body
	// it takes one arguement, maxMemory
	err := r.ParseMultipartForm(maxMemory)
	if err != nil {
		msg := "Error parsing multipart form data: " + err.Error()
		helpers.SendJSONError(w, msg, http.StatusInternalServerError)
	}

	// after parsing, now functions like r.FormValue() and r.FormFile() can be used to access the form data
	file, header, err := r.FormFile("file")
	if err != nil {
		msg := "Error accessing video file: " + err.Error()
		helpers.SendJSONError(w, msg, http.StatusInternalServerError)
		return
	}

	defer file.Close()

	contentType, err := helpers.VerifyContentType(file)
	if err != nil {
		msg := "Error verifying file content type: " + err.Error()
		helpers.SendJSONError(w, msg, http.StatusInternalServerError)
		return
	}

	log.Printf("Content-Type: %s", contentType)

	// 🛑 SECURITY CHECK: Only allow videos for now
	if !strings.HasPrefix(contentType, "video/") { // intentional error throwing
		msg := "File is not a valid video (" + contentType + ")"
		helpers.SendJSONError(w, msg, http.StatusUnsupportedMediaType)
		return
	}

	objectName := uuid.New().String() + ".mp4"
	fileSizeBytes := header.Size

	log.Printf("File size: %d", fileSizeBytes)

	info, err := h.StorageClient.PutObject(
		context.Background(),
		h.Env.MinioBucketName,
		objectName,
		file,
		fileSizeBytes,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)

	if err != nil {
		msg := "Error uploading video file: " + err.Error()
		helpers.SendJSONError(w, msg, http.StatusInternalServerError)
		return
	}

	log.Printf("Upload size: %d", info.Size)

	videoID := strings.Split(objectName, ".")[0]

	payload := job.Payload{
		"video_id": videoID,
	}

	j := job.NewJob(
		job.JobTypeTranscoding,
		payload,
		job.JobStatusPending,
	)

	h.Pool.Submit(j)

	status := h.Pool.GetJobStatus(j.ID)

	log.Printf("Job queued: %s\n", videoID)

	filename := header.Filename
	statusURL := fmt.Sprintf("%s/api/status/%s", h.Env.BaseUrl, j.ID)
	streamURL := fmt.Sprintf("%s/api/stream/%s/index.m3u8", h.Env.BaseUrl, videoID)
	videoURL := fmt.Sprintf("%s/video/%s", h.Env.ClientBaseUrl, videoID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(UploadResponse{
		Success: true,
		Data: UploadData{
			VideoID: videoID,
			JobID:   j.ID,
			Status:  status,
			Metadata: FileMetaData{
				Filename:    filename,
				SizeBytes:   fileSizeBytes,
				ContentType: contentType,
			},
			Links: ActionLinks{
				StatusURL: statusURL,
				StreamURL: streamURL,
				VideoURL:  videoURL,
			},
			CreatedAt: j.CreatedAt,
		},
	})
}
