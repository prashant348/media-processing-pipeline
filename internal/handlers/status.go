package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"media_processing_pipeline/internal/helpers"
	"media_processing_pipeline/internal/job"
	"net/http"
)

type StatusResponse struct {
	JobID  string        `json:"job_id"`
	Status job.JobStatus `json:"status"`
}

func (h *Handler) StatusHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	jobID := r.PathValue("job_id")
	
	log.Printf("Status requested for job: %s", jobID)
	
	status := h.Pool.GetJobStatus(jobID)

	if status == "unknown" {
		msg := fmt.Sprintf("Job %s not found", jobID)
		helpers.SendJSONError(w, msg, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(StatusResponse{
		JobID: jobID,
		Status: status,
	})
}
