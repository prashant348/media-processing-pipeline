package worker

import (
	"log"
	"media_processing_pipeline/internal/jobs"
	"time"
)


func InitWorkers(numberOfWorkers int, queue chan jobs.Job) {
	for i := range numberOfWorkers {
		go StartWorker(i + 1, queue)
	}
}

func StartWorker(id int, queue chan jobs.Job) {
	go func() {
		log.Printf("Worker %d started\n", id)

		for job := range queue {
			processJob(id, job)
		}

	}()
}

func processJob(workerID int, job jobs.Job) {
	log.Printf("Worker %d picked job: %s\n", workerID, job.VideoID)
	time.Sleep(5 * time.Second)
	log.Printf("Worker %d finished processing: %s\n", workerID, job.VideoID)
}
