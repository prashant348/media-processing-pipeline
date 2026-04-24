package queue

import (
	"errors"
	"media_processing_pipeline/internal/job"
)

var ErrQueueClosed = errors.New("Queue is closed.")

type Queue struct {
	Queue    chan *job.Job
	Capacity int
	IsClosed bool
}

func NewQueue(capacity int) *Queue {
	return &Queue{
		Queue:    make(chan *job.Job, capacity),
		Capacity: capacity,
		IsClosed: false,
	}
}

func (q *Queue) Enqueue(j *job.Job) error {
	if q.IsClosed {
		return ErrQueueClosed
	}

	q.Queue <- j

	return nil
}


func (q *Queue) Close() {
	if !q.IsClosed {
		close(q.Queue)
		q.IsClosed = true
	}
}
