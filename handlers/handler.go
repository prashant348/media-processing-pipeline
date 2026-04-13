package handlers

import "media_processing_pipeline/internal/jobs"

type Handler struct {
	Queue chan jobs.Job
}
