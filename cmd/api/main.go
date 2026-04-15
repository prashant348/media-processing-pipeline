package main

import (
	"fmt"
	"log"
	"media_processing_pipeline/internal/config"
	"media_processing_pipeline/internal/handlers"
	"media_processing_pipeline/internal/queue"
	"media_processing_pipeline/internal/storage"
	"media_processing_pipeline/internal/worker"
	"net/http"
	"sync"
)

func main() {
	// load environment variables
	if err := config.InitEnv(); err != nil {
		log.Fatal(err)
	}
	// access env vars
	env := config.LoadEnv()

	// initialize minio clint
	storage.InitMinIO(env)
	// acccess the initialized minio client
	client := storage.MinioClient
	// initialize and access queue
	queue := queue.InitQueue(100)
	// create wait group
	wg := &sync.WaitGroup{}
	// create worker pool
	pool := &worker.WorkerPool{
		Queue:         queue,
		WorkerCount:   3,
		Env:           env,
		StorageClient: client,
		WaitGroup:     wg,
	}

	// start the worker pool
	pool.Start()

	// create handler
	handler := &handlers.Handler{
		Pool:          pool,
		Env:           env,
		StorageClient: client,
	}

	// initialize minio
	fs := http.FileServer(http.Dir("output"))

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handler.HomeHandler)
	mux.HandleFunc("GET /video/", handler.VideoHandler)
	mux.Handle("GET /api/stream/", http.StripPrefix("/api/stream/", fs))
	// pass minio client and bucket name as Dependency injection to upload handler
	mux.HandleFunc("POST /api/upload", handler.UploadHandler())

	fmt.Printf("Server is running http://localhost:8080\n")
	log.Fatal(http.ListenAndServe(":8080", mux))

}
