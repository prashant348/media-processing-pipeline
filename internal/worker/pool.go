package worker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"media_processing_pipeline/internal/config"
	"media_processing_pipeline/internal/job"

	"media_processing_pipeline/internal/queue"
	"media_processing_pipeline/internal/storage"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/minio/minio-go/v7"
)

type WorkerPool struct {
	Queue         *queue.Queue
	WorkerCount   int
	Env           *config.Env
	StorageClient storage.ObjectStore
	WaitGroup     *sync.WaitGroup
	JobStore      *job.JobStore
}

type WorkerPoolInterface interface {
	Submit(j *job.Job)
	GetQueue() *queue.Queue
	GetJobStatus(jobID string) job.JobStatus
}

func NewWorkerPool(
	queue *queue.Queue,
	workerCount int,
	env *config.Env,
	storageClient storage.ObjectStore,
	waitGroup *sync.WaitGroup,
	jobStore *job.JobStore,
) (*WorkerPool, error) {

	if workerCount < 1 {
		return nil, errors.New("Worker count can't be zero or negative!")
	}

	return &WorkerPool{
		Queue: queue,
		WorkerCount: workerCount,
		Env: env,
		StorageClient: storageClient,
		WaitGroup: waitGroup,
		JobStore: jobStore,
	}, nil
}

func (wp *WorkerPool) GetQueue() *queue.Queue {
	return wp.Queue
}

func (wp *WorkerPool) GetJobStatus(jobID string) job.JobStatus {
	return wp.JobStore.GetStatus(jobID)
}

func (wp *WorkerPool) Start() {
	for i := range wp.WorkerCount {
		go wp.worker(i + 1)
	}
}

func (wp *WorkerPool) worker(id int) {
	log.Printf("Worker %d initialized", id)
	for j := range wp.Queue.Queue {
		func() {
			defer wp.WaitGroup.Done()

			wp.JobStore.UpdateStatus(j.ID, job.JobStatusProcessing)

			err, _ := wp.TranscodingJob(id, j, wp.Env)

			if err != nil {
				log.Printf("Worker %d: Error processing job %s: %v", id, j.ID, err)
				wp.JobStore.SetError(j.ID, err.Error())
				wp.JobStore.UpdateStatus(j.ID, job.JobStatusFailed)
			} else {
				log.Printf("Worker %d: Job %s completed", id, j.ID)
				wp.JobStore.UpdateStatus(j.ID, job.JobStatusCompleted)
			}
		}()
	}
}

func (wp *WorkerPool) Submit(j *job.Job) {
	wp.WaitGroup.Add(1)
	wp.JobStore.Create(j)

	err := wp.Queue.Enqueue(j)
	if err != nil {
		log.Printf("Submit failed: %s", err)
		wp.WaitGroup.Done()
	}
}

func (wp *WorkerPool) TranscodingJob(
	workerID int,
	j *job.Job,
	env *config.Env,
) (error, bool) {
	log.Printf("Worker %d picked job: %s\n", workerID, j.ID)

	inputPath := fmt.Sprintf("tmp/%s.mp4", j.Payload["video_id"])

	inputDir := "tmp"
	outputDir := fmt.Sprintf("output/%s", j.Payload["video_id"])

	os.MkdirAll(inputDir, os.ModePerm)
	os.MkdirAll(outputDir, os.ModePerm)

	objectName := j.Payload["video_id"] + ".mp4" 

	obj, err := wp.StorageClient.GetObject(
		context.Background(),
		env.MinioBucketName,
		objectName,
		minio.GetObjectOptions{},
	)

	if err != nil {
		log.Println("TranscodingJob: Error getting object: ", err)
		return err, false
	}
	defer obj.Close()

	file, err := os.Create(inputPath)
	if err != nil {
		log.Println("TranscodingJob: Error creating file: ", err)
		return err, false
	}
	defer file.Close()

	_, err = io.Copy(file, obj)
	if err != nil {
		log.Println("TranscodingJob: Error copying object: ", err)
		return err, false
	}

	cmd := exec.Command(
		"ffmpeg",        // ffmpeg cmd
		"-i", inputPath, // input file path
		"-codec", "copy", // direct copy streams
		"-start_number", "0", // segmentation numbers starts with 0
		"-hls_time", "5", //
		"-hls_list_size", "0", //
		"-f", "hls", // output format
		filepath.Join(outputDir, "index.m3u8"),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Println("Starting FFmpeg...")

	if err := cmd.Run(); err != nil {
		log.Printf("TranscodingJob: FFmpeg failed for %s: %v\n", j.Payload["video_id"], err)
		return err, false
	}

	log.Printf("Worker %d finished processing: %s\n", workerID, j.Payload["video_id"])

	return nil, true
}
