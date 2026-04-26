package handlers

import (
	"encoding/json"
	"fmt"
	"media_processing_pipeline/internal/config"
	"media_processing_pipeline/internal/handlers"
	"media_processing_pipeline/internal/helpers"
	"media_processing_pipeline/internal/job"
	"media_processing_pipeline/internal/queue"
	"media_processing_pipeline/internal/worker"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestSetAndGetStatus(t *testing.T) {

	jobStore := job.NewJobStore()

	j := &job.Job{
		ID:     "random-job-id",
		Status: job.JobStatusPending,
	}

	jobStore.Create(j)

	jobStore.UpdateStatus(j.ID, job.JobStatusPending)

	jobStore.UpdateStatus(j.ID, job.JobStatusProcessing)

	jobStore.UpdateStatus(j.ID, job.JobStatusCompleted)

	status := jobStore.GetStatus(j.ID)

	t.Logf("status: %s", status)

	if status != job.JobStatusCompleted {
		t.Fatalf("expected status to be: %s, got %s", job.JobStatusCompleted, status)
	}
}

func TestGetStatusOfUnknownJob(t *testing.T) {

	jobStore := job.NewJobStore()

	jobID := "random-job-id"
	// do not create job
	// jobStore.Create(j)

	status := jobStore.GetStatus(jobID)

	t.Logf("status: %s", status)

	expected_status := "unknown"
	if status != "unknown" {
		t.Fatalf("expected %s, got: %s", expected_status, status)
	}
}

func TestStatusHandlerWithValidJob(t *testing.T) {

	jobStore := job.NewJobStore()

	j := &job.Job{
		ID: "random-job-id",
	}

	jobStore.Create(j)

	jobStore.UpdateStatus(j.ID, job.JobStatusCompleted)

	status := jobStore.GetStatus(j.ID)

	t.Logf("status: %s", status)

	req := httptest.NewRequest("GET", "/api/status/{job_id}", nil)

	req.SetPathValue("job_id", j.ID)

	rec := httptest.NewRecorder()

	mockClient := &MockStore{}

	queue, _ := queue.NewQueue(10)

	pool := &worker.WorkerPool{
		Queue:         queue,
		WorkerCount:   3,
		Env:           &config.Env{},
		StorageClient: mockClient,
		WaitGroup:     &sync.WaitGroup{},
		JobStore:      jobStore,
	}

	pool.Start()

	handler := &handlers.Handler{
		Pool:          pool,
		StorageClient: mockClient,
		Env:           &config.Env{},
	}

	handler.StatusHandler(rec, req)

	responseBody := rec.Body.Bytes()

	t.Logf("response body: %s", string(responseBody))

	jsonResponse := &handlers.StatusResponse{}

	err := json.Unmarshal(responseBody, jsonResponse)
	if err != nil {
		t.Fatalf("Failed to pass json: %s", err)
	}

	if jsonResponse.Status != "completed" {
		t.Fatalf("Expected status to be: %s, got: %s", job.JobStatusCompleted, jsonResponse.Status)
	}

	if jsonResponse.JobID != j.ID {
		t.Fatalf("Expected jobID to be: %s, got: %s", j.ID, jsonResponse.JobID)
	}
}

func TestStatusHandlerWithInvalidJob(t *testing.T) {

	jobStore := job.NewJobStore()

	j := &job.Job{
		ID: "random-job-id",
	}
	// do not create job
	// jobStore.Create(j)

	// jobStore.UpdateStatus(j.ID, job.JobStatusCompleted)

	req := httptest.NewRequest("GET", "/api/status/{job_id}", nil)

	req.SetPathValue("job_id", j.ID)

	rec := httptest.NewRecorder()

	mockClient := &MockStore{}

	queue, _ := queue.NewQueue(10)

	pool := &worker.WorkerPool{
		Queue:         queue,
		WorkerCount:   3,
		Env:           &config.Env{},
		StorageClient: mockClient,
		WaitGroup:     &sync.WaitGroup{},
		JobStore:      jobStore,
	}

	pool.Start()

	handler := &handlers.Handler{
		Pool:          pool,
		StorageClient: mockClient,
		Env:           &config.Env{},
	}

	handler.StatusHandler(rec, req)

	responseBody := rec.Body.Bytes()

	jsonResponse := &helpers.ErrorResponse{}

	err := json.Unmarshal(responseBody, jsonResponse)
	if err != nil {
		t.Fatalf("Error parsing json: %s", err)
	}

	if rec.Code != 404 {
		t.Fatalf("Expected status code to be: %d, got: %d", 404, rec.Code)
	}

	msg := fmt.Sprintf("Job %s not found", j.ID)
	if jsonResponse.Message != msg {
		t.Fatalf("Expected msg to be: %s, got: %s", msg, jsonResponse.Message)
	}

}

func TestStatusConcurrentAccess(t *testing.T) {
	// reset global test
	jobStore := job.NewJobStore()

	j := &job.Job{
		ID: "random-job-id",
	}

	jobStore.Create(j)

	var wg sync.WaitGroup

	numRoutines := 100

	statuses := []job.JobStatus{
		job.JobStatusPending,
		job.JobStatusProcessing,
		job.JobStatusCompleted,
		job.JobStatusFailed,
	}

	for i := 0; i < numRoutines; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			status := statuses[i%len(statuses)]

			jobStore.UpdateStatus(j.ID, status)

			_ = jobStore.GetStatus(j.ID)
		}(i)
	}

	wg.Wait()

	finalStatus := jobStore.GetStatus(j.ID)
	if finalStatus == "" {
		t.Errorf("Expected some status, go empty")
	}
}
