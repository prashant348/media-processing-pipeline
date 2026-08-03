package handlers

import (
	"encoding/json"
	"fmt"
	"media_processing_pipeline/internal/helpers"
	"media_processing_pipeline/internal/job"
	"net/http"
	"time"
)

type JobResponse struct {
	Success bool    `json:"success"`
	Data    JobData `json:"data"`
}

type JobData struct {
	JobID      string        `json:"job_id"`
	JobType    job.JobType   `json:"job_type"`
	Payload    job.Payload   `json:"payload"`
	Status     job.JobStatus `json:"status"`
	CreatedAt  time.Time     `json:"created_at"`
	StartedAt  *time.Time    `json:"started_at"`
	FinishedAt *time.Time    `json:"finished_at"`
	Error      string        `json:"error"`
}

func (h *Handler) JobHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	jobId := r.PathValue("job_id")

	j, ok := h.Pool.GetJob(jobId)

	if !ok {
		msg := fmt.Sprintf("Job %s not found", jobId)
		helpers.SendJSONError(w, msg, http.StatusNotFound)
		return
	}

	var startedAtPtr, finishedAtPtr *time.Time

	if !j.StartedAt.IsZero() {
		startedAtPtr = &j.StartedAt
	}
	if !j.FinishedAt.IsZero() {
		finishedAtPtr = &j.FinishedAt
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(JobResponse{
		Success: true,
		Data: JobData{
			JobID:      j.ID,
			JobType:    j.Type,
			Payload:    j.Payload,
			Status:     j.Status,
			CreatedAt:  j.CreatedAt,
			StartedAt:  startedAtPtr,
			FinishedAt: finishedAtPtr,
			Error:      j.LastError,
		},
	})
}
