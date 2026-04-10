package main

import (
	"fmt"
	"log"
	"media_processing_pipeline/config"
	"media_processing_pipeline/handlers"
	"media_processing_pipeline/storage"
	"net/http"
)

func main() {
	// load environment variables
	if err := config.InitEnv(); err != nil {
		log.Fatal(err)
	}

	env := config.LoadEnv()

	// initialize minio 
	storage.InitMinIO(env)
	// create minio client
	client := storage.MinioClient
	
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handlers.HomeHandler)
	// pass minio client and bucket name as Dependency injection to upload handler
	mux.HandleFunc("POST /upload", handlers.UploadHandler(client, env.MinioBucketName))
	
	fmt.Printf("Server is running http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

}