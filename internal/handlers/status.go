package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"media_processing_pipeline/internal/jobs"
	"net/http"
)

func (h *Handler) StatusHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	jobID := r.PathValue("job_id")
	
	log.Printf("Status requested for job: %s", jobID)
	
	status := jobs.GetStatus(jobID)

	if status == "" {
		http.Error(w, fmt.Sprintf("Job %s not found", jobID), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"job_id": jobID,
		"status": string(status),
	})
}
