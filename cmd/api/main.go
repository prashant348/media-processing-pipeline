package main

import (
	"fmt"
	"log"
	"media_processing_pipeline/internal/config"
	"media_processing_pipeline/internal/handlers"
	"media_processing_pipeline/internal/job"
	"media_processing_pipeline/internal/queue"
	"media_processing_pipeline/internal/storage"
	"media_processing_pipeline/internal/worker"
	"net/http"
	"sync"

	"github.com/rs/cors"
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
	// create job store
	jobStore := job.NewJobStore()
	// create jobs queue
	queue := queue.NewQueue(100)
	// create wait group
	wg := &sync.WaitGroup{}
	// create worker pool
	pool := &worker.WorkerPool{
		Queue:         queue,
		WorkerCount:   3,
		Env:           env,
		StorageClient: client,
		WaitGroup:     wg,
		JobStore: jobStore,
	}
	// start the worker pool
	pool.Start()

	// create handler
	handler := &handlers.Handler{
		Pool:          pool,
		Env:           env,
		StorageClient: client,
	}

	mux := http.NewServeMux()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowCredentials: true,
		AllowedMethods: []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		Debug: true, // development ke time logs dikhayega
	})

	mux.HandleFunc("GET /", handler.HomeHandler)
	mux.HandleFunc("GET /api/stream/{video_id}/index.m3u8", handler.StreamHandler)
	mux.HandleFunc("POST /api/upload", handler.UploadHandler())
	mux.HandleFunc("GET /api/status/{job_id}", handler.StatusHandler)


	fmt.Printf("Server is running http://localhost:8080\n")
	log.Fatal(http.ListenAndServe(":8080", c.Handler(mux)))

}
