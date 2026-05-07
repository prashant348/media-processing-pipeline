package unit

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"media_processing_pipeline/internal/config"
	"media_processing_pipeline/internal/handlers"
	"media_processing_pipeline/internal/job"
	"media_processing_pipeline/internal/queue"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	// "net/textproto"
	"os"
	"path/filepath"
	"testing"

	minio "github.com/minio/minio-go/v7"
)

// MockStore is a mock implementation of the MinIO client for testing
type MockStore struct{}

func (m *MockStore) PutObject(
	ctx context.Context,
	bucketName string,
	objectName string,
	reader io.Reader,
	size int64,
	opts minio.PutObjectOptions,
) (minio.UploadInfo, error) {
	return minio.UploadInfo{}, nil // just
}

func (m *MockStore) GetObject(
	ctx context.Context,
	bucketName string,
	objectName string,
	opts minio.GetObjectOptions,
) (*minio.Object, error) {
	return nil, nil
}

type MockWorkerPool struct {
	JobStore *job.JobStore
}

func (mwp *MockWorkerPool) Start() {}

func (mwp *MockWorkerPool) Submit(job *job.Job)                     {}
func (mwp *MockWorkerPool) GetQueue() *queue.Queue                  { return nil }
func (mwp *MockWorkerPool) GetJobStatus(jobID string) job.JobStatus { return job.JobStatusPending }
func (mwp *MockWorkerPool) GetJob(jobID string) (*job.Job, bool)    { return mwp.JobStore.Get(jobID) }
func (mwp *MockWorkerPool) GetJobStore() *job.JobStore              { return nil }

func TestUploadHandlerWithRealFile(t *testing.T) {
	// this is the pointer to the empty buffer created
	body := &bytes.Buffer{}
	// create the formatter/writer on/for that empty buffer
	writer := multipart.NewWriter(body)
	// create a part
	part, err := writer.CreateFormFile("file", "test_video.mp4")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("body: %s", body)
	// get the path to the file
	filepath := filepath.Join("..", "testdata", "test_video.mp4")
	// read the file
	file, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}
	// write dummy video file bytes to the part
	part.Write(file)
	// finish the formatting by
	writer.Close()
	// create a fake request
	req := httptest.NewRequest("POST", "/api/upload", body)
	// tell the server this is a "multipart" upload by setting "Content-Type" field of request header
	req.Header.Set("Content-Type", writer.FormDataContentType())
	// create a "Recorder" (the spy)
	// In Go, this 'rec' acts like a browser waiting for a response from server
	rec := httptest.NewRecorder()
	// create a fake minio client
	mockClient := &MockStore{}
	// create a fake worker pool
	mockPool := &MockWorkerPool{}

	mockPool.Start()

	handler := &handlers.Handler{
		Pool:          mockPool,
		StorageClient: mockClient,
		Env: &config.Env{
			MinioBucketName: "videos",
		},
	}

	uploadHandler := handler.UploadHandler
	uploadHandler(rec, req)

	// assesrtions
	if rec.Code != http.StatusOK {
		t.Errorf("Err: %s", rec.Body.String())
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	responseBody := rec.Body.Bytes()
	uploadResponse := &handlers.UploadResponse{}

	err = json.Unmarshal(responseBody, uploadResponse)
	if err != nil {
		t.Fatalf("Failed to parse json: %s", err)
	}

	if uploadResponse.Data.VideoID == "" {
		t.Errorf("Expected video ID, got empty string")
	}
}

func TestUploadHandlerWithDummyFile(t *testing.T) {
	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "test_video.mp4")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("body: %s", body)

	filepath := filepath.Join("..", "testdata", "dummy_video.mp4")
	// read the file
	file, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

	part.Write(file)

	writer.Close()

	req := httptest.NewRequest("POST", "/api/upload", body)

	req.Header.Set("Content-Type", writer.FormDataContentType())

	rec := httptest.NewRecorder()

	mockClient := &MockStore{}

	mockPool := &MockWorkerPool{}

	mockPool.Start()

	handler := &handlers.Handler{
		Pool:          mockPool,
		StorageClient: mockClient,
		Env: &config.Env{
			MinioBucketName: "videos",
		},
	}

	uploadHandler := handler.UploadHandler
	uploadHandler(rec, req)

	// assesrtions
	if rec.Code != http.StatusUnsupportedMediaType {
		t.Errorf("Err: %s", rec.Body.String())
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	responseBody := rec.Body.Bytes()
	uploadResponse := &handlers.UploadResponse{}

	err = json.Unmarshal(responseBody, uploadResponse)
	if err != nil {
		t.Fatalf("Failed to parse json: %s", err)
	}

	if uploadResponse.Data.VideoID != "" {
		t.Errorf("Expected empty video ID, got non-empty string")
	}
}
