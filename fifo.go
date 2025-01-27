package schedule

import (
	"context"
	"sync"
)

type FIFO struct {
	jobs  []Job
	done  chan struct{}
	errCh chan error
	mu    sync.RWMutex
}

var _ Scheduler = (*FIFO)(nil)

func NewFIFO() *FIFO {
	return &FIFO{
		jobs:  make([]Job, 0),
		done:  make(chan struct{}),
		errCh: make(chan error),
	}
}

func (f *FIFO) Add(job Job) {
	f.jobs = append(f.jobs, job)
}

func (l *FIFO) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.jobs)
}

func (f *FIFO) Run(ctx context.Context) {
	if ctx == nil {
		ctx = context.TODO()
	}

	go func() {
		defer close(f.done)

		var err error
		for {
			f.mu.Lock()
			if len(f.jobs) == 0 {
				f.mu.Unlock()
				break
			}

			job := f.jobs[0]
			f.jobs = f.jobs[1:]
			f.mu.Unlock()

			select {
			case <-ctx.Done():
				if err = ctx.Err(); err != nil {
					f.errCh <- err
				}
				return
			default:
				if err = job.Done(ctx); err != nil {
					f.errCh <- err
					return
				}
			}
		}
	}()
}

func (f *FIFO) Done() <-chan struct{} {
	return f.done
}

func (f *FIFO) Errors() <-chan error {
	return f.errCh
}
