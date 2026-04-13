package main

import (
	"fmt"
	"log"
	"media_processing_pipeline/config"
	"media_processing_pipeline/handlers"
	"media_processing_pipeline/internal/queue"
	"media_processing_pipeline/internal/worker"
	"media_processing_pipeline/storage"
	"net/http"
)

func main() {
	// load environment variables
	if err := config.InitEnv(); err != nil {
		log.Fatal(err)
	}

	env := config.LoadEnv()

	queue := queue.InitQueue(100)
	worker.InitWorkers(3, queue)

	handler := &handlers.Handler{
		Queue: queue,
	}

	// initialize minio 
	storage.InitMinIO(env)
	// create minio client
	client := storage.MinioClient

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handlers.HomeHandler)
	// pass minio client and bucket name as Dependency injection to upload handler
	mux.HandleFunc("POST /api/upload", handler.UploadHandler(client, env.MinioBucketName))

	fmt.Printf("Server is running http://localhost:8080\n")
	log.Fatal(http.ListenAndServe(":8080", mux))

}