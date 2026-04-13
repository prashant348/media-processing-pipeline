package queue

import (
	"log"
	"media_processing_pipeline/internal/jobs"
)

var JobQueue chan jobs.Job

func InitQueue(bufferSize int) {
	JobQueue = make(chan jobs.Job, bufferSize)
	log.Println("Job queue initialized")
}

