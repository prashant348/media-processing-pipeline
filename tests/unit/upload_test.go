package unit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"media_processing_pipeline/internal/config"
	"media_processing_pipeline/internal/handlers"
	"media_processing_pipeline/internal/job"
	"media_processing_pipeline/internal/queue"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
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


type MockWorkerPool struct{
	JobStore *job.JobStore
}

func (mwp *MockWorkerPool) Start() {}

func (mwp *MockWorkerPool) Submit(job *job.Job) {}
func (mwp *MockWorkerPool) GetQueue() *queue.Queue { return nil }
func (mwp *MockWorkerPool) GetJobStatus(jobID string) job.JobStatus { return job.JobStatusPending }
func (mwp *MockWorkerPool) GetJob(jobID string) (*job.Job, bool) { return mwp.JobStore.Get(jobID) }

func TestUploadHandler(t *testing.T) {

	// this is the pointer to the empty buffer created
	body := &bytes.Buffer{}
	// create the formatter/writer on/for that empty buffer
	writer := multipart.NewWriter(body)
	// create a fake file inside that body
	// this mimics a user selecting a file in their browser
	// formats the empty buffer by creating a form data header with the provided fields
	// this uses up a few bytes of that buffer to write the things like Content-Disposition: form-data; name="file"; filename="test_video.mp4"
	// now the buffer is no longer empty because it contains "metadata" of the file
	part, err := writer.CreateFormFile("file", "test_video.mp4")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("----------DEBUG: What is inside the buffer now?----------------\n")
	fmt.Println(*body)            // raw bytes
	fmt.Println((*body).String()) // human readable. It will contain header (metadata) of the form data (file)
	fmt.Printf("----------------------------------------------------------------\n")
	// write some dummy bytes
	part.Write([]byte("this is fake video data"))
	// finish the formatting by
	writer.Close()
	fmt.Printf("----------DEBUG: What is inside the buffer now?----------------\n")
	fmt.Println((*body).String())
	fmt.Printf("----------------------------------------------------------------\n")
	// create a fake request
	req := httptest.NewRequest("POST", "/upload", body)
	// tell the server this is a "multipart" upload by setting "Content-Type" field of request header
	req.Header.Set("Content-Type", writer.FormDataContentType())
	// create a "Recorder" (the spy)
	// In Go, this 'rec' acts like a browser waiting for a response from server
	rec := httptest.NewRecorder()
	// create a fake minio client
	mockClient := &MockStore{}
	// run the handler by passing DIs minio client and bucket name
	// queue := queue.InitQueue(10)


	mockPool := &MockWorkerPool{
	}

	mockPool.Start()

	handler := &handlers.Handler{
		Pool:          mockPool,
		StorageClient: mockClient,
		Env: &config.Env{
			MinioBucketName: "videos",
		},
	}

	uploadHandler := handler.UploadHandler()
	uploadHandler.ServeHTTP(rec, req)

	// assesrtions
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	responseBody := rec.Body.Bytes()
	jsonResponse := &handlers.UploadResponse{}
	
	err = json.Unmarshal(responseBody, jsonResponse)
	if err != nil {
		t.Fatalf("Failed to parse json: %s", err)
	}
 
	if jsonResponse.VideoID == "" {
		t.Errorf("Expected video ID, got empty string")
	}
}
