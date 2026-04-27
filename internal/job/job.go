package job

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Job struct {
	ID        string    `json:"id"`
	Type      JobType   `json:"type"`
	Payload   Payload   `json:"payload"`
	Status    JobStatus `json:"status"`
	// job lifecycle timings
	CreatedAt time.Time `json:"created_at"`
	StartedAt time.Time `json:"started_at"`
	FinishedAt time.Time `json:"updated_at"`
	// error handling
	LastError string    `json:"last_error"`
}

type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
)

type JobType string

const (
	JobTypeTranscoding JobType = "transcoding"
)

type Payload map[string]string

type JobStore struct {
	mu   sync.Mutex
	Jobs map[string]*Job
}

func NewJob(
	jobType JobType,
	payload Payload,
	status JobStatus,
) *Job {
	return &Job{
		ID: uuid.New().String(),
		Type: jobType,
		Payload: payload,
		Status: status,
		LastError: "",
		CreatedAt: time.Now(),
		StartedAt: time.Time{},
		FinishedAt: time.Time{},
	}
}

func NewJobStore() *JobStore {
	return &JobStore{Jobs: make(map[string]*Job)}
}

func (js *JobStore) Get(jobID string) (*Job, bool) {
	js.mu.Lock()
	defer js.mu.Unlock()

	j, ok := js.Jobs[jobID]
	return j, ok
}

func (js *JobStore) Create(job *Job) {
	js.mu.Lock()
	defer js.mu.Unlock()

	js.Jobs[job.ID] = job
}

func (js *JobStore) UpdateStatus(
	jobID string,
	status JobStatus,
) {
	js.mu.Lock()
	defer js.mu.Unlock()

	js.Jobs[jobID].Status = status
}

func (js *JobStore) UpdateStartedAt(jobID string) {
	js.mu.Lock()
	defer js.mu.Unlock()

	js.Jobs[jobID].StartedAt = time.Now()
}

func (js *JobStore) UpdateFinishedAt(jobID string) {
	js.mu.Lock()
	defer js.mu.Unlock()

	js.Jobs[jobID].FinishedAt = time.Now()
}

func (js *JobStore) GetStatus(jobID string) JobStatus {
	js.mu.Lock()
	defer js.mu.Unlock()

	if j, ok := js.Jobs[jobID]; ok {
		return j.Status
	}

	return "unknown"
}

func (js *JobStore) SetError(
	jobID string,
	err string,
) {
	js.mu.Lock()
	defer js.mu.Unlock()

	js.Jobs[jobID].LastError = err
}
