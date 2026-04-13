package worker

import (
	"log"
	"media_processing_pipeline/internal/jobs"
	"media_processing_pipeline/internal/queue"
	"time"
)


func InitWorkers(numberOfWorkers int) {
	for i := range numberOfWorkers {
		StartWorker(i + 1)
	}
}

func StartWorker(id int) {
	go func() {
		log.Printf("Worker %d started\n", id)

		for job := range queue.JobQueue {
			processJob(id, job)
		}

	}()
}

func processJob(workerID int, job jobs.Job) {
	log.Printf("Worker %d picked job: %s\n", workerID, job.VideoID)
	time.Sleep(5 * time.Second)
	log.Printf("Worker %d finished processing: %s\n", workerID, job.VideoID)
}
