package schedule

import (
	"context"
	"sync"
)

// PriorityJob represents a job with an associated priority level.
// It embeds the Job interface and adds a method to retrieve the priority.
type PriorityJob interface {
	Job

	// Priority returns the priority level of the job.
	Priority() int
}

// Internal implementation of a PriorityJob.
type priorityJobImpl struct {
	job      Job
	priority int
}

var _ Job = (*priorityJobImpl)(nil)

func NewPriorityJob(ctx context.Context, handler func() error, prio int) PriorityJob {
	return &priorityJobImpl{
		job:      NewJob(ctx, handler),
		priority: prio,
	}
}

func NewPriorityJobFrom(job Job, prio int) PriorityJob {
	return &priorityJobImpl{
		job:      job,
		priority: prio,
	}
}

func (pj *priorityJobImpl) Done() error {
	return pj.job.Done()
}

func (pj *priorityJobImpl) Context() context.Context {
	return pj.job.Context()
}

func (pj *priorityJobImpl) Priority() int {
	return pj.priority
}

// PriorityQueue
type PriorityQueue struct {
	jobs  []Job
	done  chan struct{}
	errCh chan error
	mu    sync.RWMutex
}

var _ Scheduler = (*PriorityQueue)(nil)

func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{
		jobs:  make([]Job, 0),
		done:  make(chan struct{}),
		errCh: make(chan error),
	}
}

func (pq *PriorityQueue) Add(job Job) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	newJob, ok := job.(PriorityJob)
	if !ok {
		pq.jobs = append(pq.jobs, NewPriorityJobFrom(job, 0))
		return
	}

	if len(pq.jobs) == 0 {
		pq.jobs = append(pq.jobs, job)
	} else {
		j := 0

		curJob := pq.jobs[j].(PriorityJob)
		for j < len(pq.jobs) && newJob.Priority() <= curJob.Priority() {
			j++

			if j < len(pq.jobs) {
				curJob = pq.jobs[j].(PriorityJob)
			}
		}

		if j == 0 {
			pq.jobs = append([]Job{newJob}, pq.jobs...)
		} else if j == len(pq.jobs) {
			pq.jobs = append(pq.jobs, newJob)
		} else {
			pq.jobs = append(pq.jobs[:j+1], pq.jobs[j:]...)
			pq.jobs[j] = newJob
		}
	}
}

func (pq *PriorityQueue) Len() int {
	pq.mu.RLock()
	defer pq.mu.RUnlock()
	return len(pq.jobs)
}

func (pq *PriorityQueue) Run(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}

	go func() {
		defer close(pq.done)

		var err error
		for {
			pq.mu.Lock()
			if len(pq.jobs) == 0 {
				pq.mu.Unlock()
				break
			}

			job := pq.jobs[0]
			pq.jobs = pq.jobs[1:]
			pq.mu.Unlock()

			select {
			case <-ctx.Done():
				if err = ctx.Err(); err != nil {
					pq.errCh <- err
				}
				return
			default:
				if err = job.Done(); err != nil {
					pq.errCh <- err
					return
				}
			}
		}
	}()
}

func (pq *PriorityQueue) Done() <-chan struct{} {
	return pq.done
}

func (pq *PriorityQueue) Errors() <-chan error {
	return pq.errCh
}
