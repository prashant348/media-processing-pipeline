package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"media_processing_pipeline/internal/config"
	"media_processing_pipeline/internal/handlers"
	"media_processing_pipeline/internal/job"

	"media_processing_pipeline/internal/queue"
	"media_processing_pipeline/internal/worker"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"

	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	miniodriver "github.com/testcontainers/testcontainers-go/modules/minio"
)


func TestUploadToProcessPipeline(t *testing.T) {

	ctx := context.Background()

	minioContainer, err := miniodriver.Run(
		ctx,
		"minio/minio:latest",
	)

	if err != nil {
		t.Fatal(err)
	}

	defer minioContainer.Terminate(ctx)

	endpoint, _ := minioContainer.ConnectionString(ctx)

	t.Logf("endpoint: %s", endpoint)

	realClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})

	if err != nil {
		t.Fatalf("failed to create minio realClient: %v", err)
	}

	realClient.MakeBucket(ctx, "videos", minio.MakeBucketOptions{})

	videoFilePath := filepath.Join("..", "testdata", "tiny_test_video.mp4")

	t.Logf("video file path: %s", videoFilePath)

	file, err := os.Open(videoFilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "tiny_test_video.mp4")

	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal(err)
	}

	writer.Close()

	req := httptest.NewRequest("POST", "/api/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	jobStore := job.NewJobStore()
	queue := queue.NewQueue(10)
	wg := &sync.WaitGroup{}
	pool := worker.NewWorkerPool(
		queue,
		3,
		&config.Env{
			MinioBucketName: "videos",
		},
		realClient,
		wg,
		jobStore,
	)

	pool.Start()

	handler := handlers.NewHandler(
		pool,
		realClient,
		&config.Env{
			MinioBucketName: "videos",
		},
	)
	
	uploadHandler := handler.UploadHandler()
	uploadHandler.ServeHTTP(rec, req)

	httpResponse := rec.Body.Bytes()

	jsonResponse := &handlers.UploadResponse{}
	
	err = json.Unmarshal(httpResponse, jsonResponse)
	if err != nil {
		t.Fatalf("Failed to parse json: %s", err)
	}

	if jsonResponse.VideoID == "" {
		t.Errorf("Expected video ID, got empty string")
	}

	videoID := jsonResponse.VideoID
	jobID := jsonResponse.JobID
	status := jsonResponse.Status

	if status != job.JobStatusPending && status != job.JobStatusProcessing {
		t.Errorf("Expected status to be %s or %s, got %s", job.JobStatusPending, job.JobStatusProcessing, status)
	}

	// define output path
	outputPath := filepath.Join("output", videoID, "index.m3u8")
	t.Log(outputPath)
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	pool.WaitGroup.Wait()

	// use output as assertion
	_, err = os.Stat(outputPath)
	if err != nil {
		t.Fatal("HLS playlist not generated")
	}

	status = jobStore.GetStatus(jobID)

	if status != job.JobStatusCompleted {
		t.Errorf("Exected status to be %s, got %s", job.JobStatusCompleted, status)
	}
}
