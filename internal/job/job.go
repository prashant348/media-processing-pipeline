package job

import (
	"sync"
	"time"
)

type Job struct {
	ID        string    `json:"id"`
	Type      JobType   `json:"type"`
	Payload   Payload   `json:"payload"`
	Status    JobStatus `json:"status"`
	LastError string    `json:"last_error"`
	CreatedAt time.Time `json:"created_at"`
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

func NewJobStore() *JobStore {
	return &JobStore{Jobs: make(map[string]*Job)}
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
