package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"media_processing_pipeline/internal/helpers"
	"media_processing_pipeline/internal/job"
	"net/http"
	"time"
)

type StatusResponse struct {
	Success bool       `json:"success"`
	Data    StatusData `json:"data"`
}

type StatusData struct {
	JobID      string        `json:"job_id"`
	Status     job.JobStatus `json:"status"`
	Payload    job.Payload   `json:"payload"`
	Timings    Timings       `json:"timings"`
	Error      string        `json:"error"`
	IsFinished bool          `json:"is_finished"`
}

type Timings struct {
	CreatedAt  time.Time `json:"created_at"`
	StartedAt  *time.Time `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
	Duration   string    `json:"duration"`
	QueueWait  string    `json:"queue_wait"`
}

func (h *Handler) StatusHandler(w http.ResponseWriter, r *http.Request) {
	jobID := r.PathValue("job_id")
	log.Printf("Status requested for job: %s", jobID)

	j, ok := h.Pool.GetJob(jobID)
	if !ok {
		helpers.SendJSONError(w, fmt.Sprintf("Job %s not found", jobID), http.StatusNotFound)
		return
	}

	var startedAtPtr, finishedAtPtr *time.Time

	if !j.StartedAt.IsZero() {
		startedAtPtr = &j.StartedAt
	}
	if !j.FinishedAt.IsZero() {
		finishedAtPtr = &j.FinishedAt
	}

	var duration, queueWait string

	// 1. Calculate Queue Wait (Jab job pick ho chuki ho)
	if !j.StartedAt.IsZero() {
		queueWait = j.StartedAt.Sub(j.CreatedAt).String()
	} else {
		queueWait = "waiting..."
	}

	// 2. Calculate Duration
	if !j.FinishedAt.IsZero() {
		// Completed or Failed
		duration = j.FinishedAt.Sub(j.StartedAt).String()
	} else if !j.StartedAt.IsZero() {
		// Currently Processing: Show elapsed time
		duration = time.Since(j.StartedAt).Round(time.Second).String()
	} else {
		duration = "not started"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(StatusResponse{
		Success: true,
		Data: StatusData{
			JobID:   jobID,
			Status:  j.Status,
			Payload: j.Payload,
			Timings: Timings{
				CreatedAt:  j.CreatedAt,
				StartedAt:  startedAtPtr,
				FinishedAt: finishedAtPtr,
				Duration:   duration,
				QueueWait:  queueWait,
			},
			Error:      j.LastError,
			IsFinished: j.Status == job.JobStatusCompleted || j.Status == job.JobStatusFailed,
		},
	})
}
