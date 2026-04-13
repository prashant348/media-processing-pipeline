package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"media_processing_pipeline/internal/queue"
	"media_processing_pipeline/internal/worker"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	miniodriver "github.com/testcontainers/testcontainers-go/modules/minio"
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

func TestUploadHandler(t *testing.T) {

	queue.InitQueue(100)
	worker.InitWorkers(3)
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
	handler := UploadHandler(mockClient, "videos")
	handler.ServeHTTP(rec, req)

	// assesrtions
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	expectedPrefix := "Uploaded"
	if !bytes.HasPrefix(rec.Body.Bytes(), []byte(expectedPrefix)) {
		t.Errorf("expected response to start with %s, got %s", expectedPrefix, rec.Body.String())
	}
}

func TestUploadHandlerIntegration(t *testing.T) {

	queue.InitQueue(100)
	worker.InitWorkers(3)

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
		Creds: credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})

	if err != nil {
		t.Fatalf("failed to create minio realClient: %v", err)
	}

	realClient.MakeBucket(ctx, "videos", minio.MakeBucketOptions{})

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "sample_video.mp4")
	part.Write([]byte("actual binary video data"))
	writer.Close()
 
	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	handler := UploadHandler(realClient, "videos")
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	expectedPrefix := "Uploaded"
	if !bytes.HasPrefix(rec.Body.Bytes(), []byte(expectedPrefix)) {
		t.Errorf("expected response to start with %s, got %s", expectedPrefix, rec.Body.String())
	} 
}