package worker

import (
	"log"
	"media_processing_pipeline/internal/jobs"
	"time"
)

type WorkerPool struct {
	Queue chan jobs.Job
	WorkerCount int
}

func (wp *WorkerPool) Start() {
	for i := range wp.WorkerCount {
		go wp.worker(i + 1)
	}
}

func (wp *WorkerPool) worker(id int) {
	for job := range wp.Queue {
		log.Printf("Worker %d started", id)
		processJob(id, job)
	}
}

func (wp *WorkerPool) Submit(job jobs.Job) {
	wp.Queue <- job
}

func processJob(workerID int, job jobs.Job) {
	log.Printf("Worker %d picked job: %s\n", workerID, job.VideoID)
	time.Sleep(5 * time.Second)
	log.Printf("Worker %d finished processing: %s\n", workerID, job.VideoID)
}
