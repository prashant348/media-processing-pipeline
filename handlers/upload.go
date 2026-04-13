package handlers

import (
	"context"
	"fmt"
	"log"
	"media_processing_pipeline/internal/jobs"
	"media_processing_pipeline/storage"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func (h *Handler) UploadHandler(client storage.ObjectStore, bucketName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// set a memory limit of 100KB
		const maxMemory = 100 << 10
		// this method is used to process/parse "multipart/form-data" MIME type request body
		// it takes one arguement, maxMemory
		err := r.ParseMultipartForm(maxMemory)
		if err != nil {
			http.Error(w, "Error parsing multipart form", 400)
			return
		}

		// after parsing, now functions like r.FormValue() and r.FormFile() can be used to access the form data
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error accessing file", 400)
			return
		}

		defer file.Close()

		objectName := uuid.New().String() + ".mp4"

		_, err = client.PutObject(
			context.Background(),
			bucketName,
			objectName,
			file,
			header.Size,
			minio.PutObjectOptions{
				ContentType: "application/octet-stream",
			},
		)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		videoID := strings.Split(objectName, ".")[0]

		job := jobs.Job{
			VideoID: videoID,
			FileKey: objectName,
		}

		h.Pool.Submit(job)

		log.Printf("Job queued: %s\n", videoID)

		fmt.Fprintf(w, "Uploaded %s", objectName)
	}

}
