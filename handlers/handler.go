package handlers

import (
	"media_processing_pipeline/internal/worker"
)

type Handler struct {
	Pool *worker.WorkerPool
}
