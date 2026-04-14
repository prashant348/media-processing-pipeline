package handlers

import (
	"media_processing_pipeline/config"
	"media_processing_pipeline/internal/worker"
	"media_processing_pipeline/storage"
)

type Handler struct {
	Pool *worker.WorkerPool
	StorageClient storage.ObjectStore
	Env *config.Env
}
