package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"media_processing_pipeline/internal/config"
	"media_processing_pipeline/internal/handlers"
	"media_processing_pipeline/internal/jobs"
	"media_processing_pipeline/internal/queue"
	"media_processing_pipeline/internal/worker"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	// "strings"
	"sync"
	"testing"

	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	miniodriver "github.com/testcontainers/testcontainers-go/modules/minio"
)

type UploadResponse struct {
	VideoID  string `json:"video_id"`
}

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

	queue := queue.InitQueue(10)
	wg := &sync.WaitGroup{}
	pool := &worker.WorkerPool{
		Queue:       queue,
		WorkerCount: 3,
		Env: &config.Env{
			MinioBucketName: "videos",
		},
		StorageClient: realClient,
		WaitGroup:     wg,
	}

	pool.Start()

	handler := &handlers.Handler{
		Pool:          pool,
		StorageClient: realClient,
		Env: &config.Env{
			MinioBucketName: "videos",
		},
	}
	uploadHandler := handler.UploadHandler()
	uploadHandler.ServeHTTP(rec, req)

	httpResponse := rec.Body.Bytes()

	jsonResponse := &UploadResponse{}
	
	err = json.Unmarshal(httpResponse, jsonResponse)
	if err != nil {
		t.Fatalf("Failed to parse json: %s", err)
	}

	if jsonResponse.VideoID == "" {
		t.Errorf("Expected video ID, got empty string")
	}

	videoID := jsonResponse.VideoID
	jobID := videoID

	status := jobs.GetStatus(jobID)

	if status != jobs.JobStatusPending && status != jobs.JobStatusProcessing {
		t.Errorf("Expected status to be %s or %s, got %s", jobs.JobStatusPending, jobs.JobStatusProcessing, status)
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

	status = jobs.GetStatus(jobID)

	if status != jobs.JobStatusCompleted {
		t.Errorf("Exected status to be %s, got %s", jobs.JobStatusCompleted, status)
	}
}
