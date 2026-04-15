package handlers

import (
	"media_processing_pipeline/internal/config"
	"media_processing_pipeline/internal/storage"
	"media_processing_pipeline/internal/worker"
)

type Handler struct {
	Pool worker.WorkerPoolInterface // pass worker pool interface instead of struct, for loose coupling
	StorageClient storage.ObjectStore
	Env *config.Env
}
