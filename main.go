package main

import (
	"fmt"
	"log"
	"media_processing_pipeline/config"
	"media_processing_pipeline/handlers"
	"media_processing_pipeline/storage"
	"net/http"
	"github.com/joho/godotenv"
)

func main() {

	// load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error in loading .env file")
	}

	env := config.LoadEnv()

	// initialize minio 
	storage.InitMinIO(env)
	
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handlers.HomeHandler)
	mux.HandleFunc("POST /upload", handlers.UploadHandler(env.MinioBucketName))
	
	fmt.Printf("Server is running http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

}