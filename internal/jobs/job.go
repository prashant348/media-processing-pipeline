package jobs

import "sync"

type Job struct {
	VideoID string
	FileKey string
}

type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
)

var (
	Jobs = make(map[string]JobStatus)
	mu   sync.RWMutex
)

func GetStatus(jobID string) JobStatus {
	mu.RLock()
	defer mu.RUnlock()
	return Jobs[jobID]
}

func SetStatus(jobID string, status JobStatus) {
	mu.Lock()
	defer mu.Unlock()
	Jobs[jobID] = status
}