package handlers

import (
	"encoding/json"
	"media_processing_pipeline/internal/job"
	"net/http"
)

type JobsResponse struct {
	Success bool       `json:"success"`
	Data    []*job.Job `json:"data"`
}

func (h *Handler) JobsHandler(w http.ResponseWriter, r *http.Request) {
	var jobs []*job.Job

	js := h.Pool.GetJobStore()

	allJobs := js.Jobs

	for _, j := range allJobs {
		jobs = append(jobs, j)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(JobsResponse{
		Success: true,
		Data:    jobs,
	})

}
