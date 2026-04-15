package worker

import (
	"context"
	"fmt"
	"io"
	"log"
	"media_processing_pipeline/internal/config"
	"media_processing_pipeline/internal/jobs"
	"media_processing_pipeline/internal/storage"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/minio/minio-go/v7"
)

type WorkerPool struct {
	Queue         chan jobs.Job
	WorkerCount   int
	Env           *config.Env
	StorageClient storage.ObjectStore
	WaitGroup     *sync.WaitGroup
}

type WorkerPoolInterface interface {
	Submit(job jobs.Job)
}

func (wp *WorkerPool) Start() {
	for i := range wp.WorkerCount {
		go wp.worker(i + 1)
	}
}

func (wp *WorkerPool) worker(id int) {
	log.Printf("Worker %d started", id)
	for job := range wp.Queue {
		wp.ProcessJob(id, job, wp.Env)
		wp.WaitGroup.Done()
	}
}

func (wp *WorkerPool) Submit(job jobs.Job) {
	wp.WaitGroup.Add(1)
	wp.Queue <- job
}

func (wp *WorkerPool) ProcessJob(workerID int, job jobs.Job, env *config.Env) {
	log.Printf("Worker %d picked job: %s\n", workerID, job.VideoID)

	inputPath := fmt.Sprintf("tmp/%s.mp4", job.VideoID)

	inputDir := "tmp"
	outputDir := fmt.Sprintf("output/%s", job.VideoID)

	os.MkdirAll(inputDir, os.ModePerm)
	os.MkdirAll(outputDir, os.ModePerm)

	obj, err := wp.StorageClient.GetObject(
		context.Background(),
		env.MinioBucketName,
		job.FileKey,
		minio.GetObjectOptions{},
	)

	if err != nil {
		log.Println("Error getting object: ", err)
		return
	}
	defer obj.Close()

	// log.Printf("From processJob func: obj: %v", obj)
	// log.Printf("size of obj: %v", unsafe.Sizeof(obj))
	stat, err := obj.Stat()
	if err != nil {
		log.Println("Stat error:", err)
	} else {
		log.Println("Object size:", stat.Size)
	}

	file, err := os.Create(inputPath)
	if err != nil {
		log.Println("Error creating file: ", err)
		return
	}
	defer file.Close()

	// log.Printf("From processJob func: file: %v", file)

	_, err = io.Copy(file, obj)
	if err != nil {
		log.Println("Error copying object: ", err)
		return
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
		log.Printf("FFmpeg failed for %s: %v\n", job.VideoID, err)
		return
	}

	log.Printf("Worker %d finished processing: %s\n", workerID, job.VideoID)
}
