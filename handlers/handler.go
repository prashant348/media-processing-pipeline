package handlers

import (
	"media_processing_pipeline/config"
	"media_processing_pipeline/internal/worker"
	"media_processing_pipeline/storage"
)

type Handler struct {
	Pool worker.WorkerPoolInterface // pass worker pool interface instead of struct, for loose coupling
	StorageClient storage.ObjectStore
	Env *config.Env
}
