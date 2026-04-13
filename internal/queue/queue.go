package queue

import (
	"log"
	"media_processing_pipeline/internal/jobs"
)


func InitQueue(bufferSize int) chan jobs.Job {
	queue := make(chan jobs.Job, bufferSize)
	log.Println("Job queue initialized")
	return queue
}

