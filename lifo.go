package schedule

import (
	"context"
	"sync"
)

type LIFO struct {
	jobs  []Job
	done  chan struct{}
	errCh chan error
	mu    sync.RWMutex
}

var _ Scheduler = (*LIFO)(nil)

func NewLIFO() *LIFO {
	return &LIFO{
		jobs:  make([]Job, 0),
		done:  make(chan struct{}),
		errCh: make(chan error),
	}
}

func (l *LIFO) Add(job Job) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.jobs = append(l.jobs, job)
}

func (l *LIFO) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.jobs)
}

func (l *LIFO) Run(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}

	go func() {
		defer close(l.done)

		var err error
		for {
			l.mu.Lock()
			if len(l.jobs) == 0 {
				l.mu.Unlock()
				break
			}
			job := l.jobs[len(l.jobs)-1]
			l.jobs = l.jobs[:len(l.jobs)-1]
			l.mu.Unlock()

			select {
			case <-ctx.Done():
				if err = ctx.Err(); err != nil {
					l.errCh <- err
				}
				return
			default:
				if err = job.Done(); err != nil {
					l.errCh <- err
					return
				}
			}
		}
	}()
}

func (l *LIFO) Done() <-chan struct{} {
	return l.done
}

func (l *LIFO) Errors() <-chan error {
	return l.errCh
}
