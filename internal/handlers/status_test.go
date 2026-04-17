package handlers

import (
	"encoding/json"
	"fmt"
	"media_processing_pipeline/internal/config"
	"media_processing_pipeline/internal/jobs"
	"net/http/httptest"
	"sync"
	"testing"
)

type JsonResponse struct {
	JobID  string `json:"job_id"`
	Status string `json:"status"`
}

func TestSetAndGetStatus(t *testing.T) {

	jobs.Jobs = make(map[string]jobs.JobStatus)

	jobID := "random-job-id"

	jobs.SetStatus(jobID, jobs.JobStatusPending)

	jobs.SetStatus(jobID, jobs.JobStatusProcessing)

	jobs.SetStatus(jobID, jobs.JobStatusCompleted)

	status := jobs.GetStatus(jobID)

	t.Logf("status: %s", status)

	if status != jobs.JobStatusCompleted {
		t.Fatalf("expected status to be: %s, got %s", jobs.JobStatusCompleted, status)
	}
}

func TestGetStatusOfUnknownJob(t *testing.T) {

	jobs.Jobs = make(map[string]jobs.JobStatus)

	jobID := "random-job-id"

	status := jobs.GetStatus(jobID)

	t.Logf("status: %s", status)

	expected_status := ""
	if status != "" {
		t.Fatalf("expected %s, got: %s",expected_status, status)
	}
}

func TestStatusHandlerWithValidJob(t *testing.T) {

	jobs.Jobs = make(map[string]jobs.JobStatus)

	jobID := "random-job-id"

	jobs.SetStatus(jobID, jobs.JobStatusCompleted)

	status := jobs.GetStatus(jobID)

	t.Logf("status: %s", status)

	req := httptest.NewRequest("GET", "/api/status/{job_id}", nil)

	req.SetPathValue("job_id", jobID)

	rec := httptest.NewRecorder()

	mockClient := &MockStore{}

	mockPool := &MockWorkerPool{}

	mockPool.Start()

	handler := &Handler{
		Pool:          mockPool,
		StorageClient: mockClient,
		Env:           &config.Env{},
	}

	handler.StatusHandler(rec, req)

	responseBody := rec.Body.Bytes()

	t.Logf("response body: %s", string(responseBody))

	jsonResponse := &JsonResponse{}

	err := json.Unmarshal(responseBody, jsonResponse)
	if err != nil {
		t.Fatalf("Failed to pass json: %s", err)
	}

	if jsonResponse.Status != "completed" {
		t.Fatalf("Expected status to be: %s, got: %s", jobs.JobStatusCompleted, jsonResponse.Status)
	}

	if jsonResponse.JobID != jobID {
		t.Fatalf("Expected jobID to be: %s, got: %s", jobID, jsonResponse.JobID)
	}

}

func TestStatusHandlerWithInvalidJob(t *testing.T) {
	
	jobs.Jobs = make(map[string]jobs.JobStatus)
	
	jobID := "invalid-job-id"

	// do not set status
	// jobs.SetStatus(jobID, jobs.JobStatusCompleted)

	req := httptest.NewRequest("GET", "/api/status/{job_id}", nil)

	req.SetPathValue("job_id", jobID)

	rec := httptest.NewRecorder()

	mockClient := &MockStore{}

	mockPool := &MockWorkerPool{}

	mockPool.Start()

	handler := &Handler{
		Pool:          mockPool,
		StorageClient: mockClient,
		Env:           &config.Env{},
	}

	handler.StatusHandler(rec, req)

	// http.Error automatically adds "\n" in the last of the msg!
	responseBody := rec.Body.String() // "Job invalid-job-id not found\n"
	// that is why we need to add "\n" at last
	// otherwise it will throw err because in go: "text" == "text\n" is false
	errTxt := fmt.Sprintf("Job %s not found\n", jobID)
	
	t.Logf("responseBody: %s", responseBody)
	t.Logf("errTxt: %s", errTxt)

	if rec.Code != 404 {
		t.Fatalf("Expected status code to be: %d, got: %d", 404, rec.Code)
	}

	if errTxt != responseBody {
		t.Fatalf("Expected response body to be: %s, got: %s", errTxt, responseBody)
	}

}

func TestStatusConcurrentAccess(t *testing.T) {
	// reset global test
	jobs.Jobs = make(map[string]jobs.JobStatus)

	jobID := "test-job"

	var wg sync.WaitGroup

	numRoutines := 100

	statuses := []jobs.JobStatus{
		jobs.JobStatusPending,
		jobs.JobStatusProcessing,
		jobs.JobStatusCompleted,
		jobs.JobStatusFailed,
	}

	for i := 0; i < numRoutines; i++ {
		wg.Add(1)

		go func (i int)  {
			defer wg.Done()

			status := statuses[i % len(statuses)]

			jobs.SetStatus(jobID, status)

			_ = jobs.GetStatus(jobID)
		}(i)
	}

	wg.Wait()

	finalStatus := jobs.GetStatus(jobID)
	if finalStatus == "" {
		t.Errorf("Expected some status, go empty")
	}
}