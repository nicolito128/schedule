package schedule

import (
	"context"
)

// Scheduler defines methods for scheduling and managing jobs.
// It allows adding new jobs, running the scheduler, and provides channels for job completion and error handling.
type Scheduler interface {
	// Add inserts a new job into the scheduler.
	Add(job Job)

	// Len gives the number of current jobs in the scheduler.
	Len() int

	// Run initiates the scheduler to start processing jobs in a new goroutine.
	Run(ctx context.Context)

	// Done provides a channel that signals when all jobs are finished.
	Done() <-chan struct{}

	// Errors provides a channel for receiving errors encountered during job execution.
	Errors() <-chan error
}
