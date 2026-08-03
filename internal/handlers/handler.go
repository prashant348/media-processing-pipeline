package handlers

import (
	"media_processing_pipeline/internal/config"
	"media_processing_pipeline/internal/storage"
	"media_processing_pipeline/internal/worker"
)

type Handler struct {
	Pool          worker.WorkerPoolInterface // pass worker pool interface instead of struct, for loose coupling
	StorageClient storage.ObjectStore
	Env           *config.Env
}

func NewHandler(
	pool worker.WorkerPoolInterface,
	storageClient storage.ObjectStore,
	env *config.Env,
) *Handler {
	return &Handler{
		Pool:          pool,
		StorageClient: storageClient,
		Env:           env,
	}
}
