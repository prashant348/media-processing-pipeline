package handlers

import (
	"context"
	"encoding/json"

	"log"
	"media_processing_pipeline/internal/helpers"
	"media_processing_pipeline/internal/job"

	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type UploadResponse struct {
	VideoID string        `json:"video_id"`
	JobID   string        `json:"job_id"`
	Status  job.JobStatus `json:"status"`
}

func (h *Handler) UploadHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// set a memory limit of 100KB
		const maxMemory = 100 << 10
		// this method is used to process/parse "multipart/form-data" MIME type request body
		// it takes one arguement, maxMemory
		err := r.ParseMultipartForm(maxMemory)
		if err != nil {
			msg := "Error parsing multipart form data: " + err.Error()
			helpers.SendJSONError(w, msg, http.StatusBadRequest)
			return
		}

		// after parsing, now functions like r.FormValue() and r.FormFile() can be used to access the form data
		file, header, err := r.FormFile("file")
		if err != nil {
			msg := "Error accessing video file: " + err.Error()
			helpers.SendJSONError(w, msg, http.StatusBadRequest)
			return
		}

		defer file.Close()

		objectName := uuid.New().String() + ".mp4"

		info, err := h.StorageClient.PutObject(
			context.Background(),
			h.Env.MinioBucketName,
			objectName,
			file,
			header.Size,
			minio.PutObjectOptions{
				ContentType: "application/octet-stream",
			},
		)

		log.Printf("Upload size: %d", info.Size)

		if err != nil {
			msg := "Error uploading video file: " + err.Error()
			helpers.SendJSONError(w, msg, http.StatusInternalServerError)
			return
		}

		videoID := strings.Split(objectName, ".")[0]

		payload := map[string]string{
			"video_id": videoID,
		}

		j := job.NewJob(
			uuid.New().String(),
			job.JobTypeTranscoding,
			payload,
			job.JobStatusPending,
		)

		h.Pool.Submit(j)

		status := h.Pool.GetJobStatus(j.ID)

		log.Printf("Job queued: %s\n", videoID)

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(UploadResponse{
			VideoID: videoID,
			JobID:   j.ID,
			Status:  status,
		})
	}
}
