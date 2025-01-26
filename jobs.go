package schedule

import "context"

// Job represents a unit of work that can be executed.
// It provides methods to mark the job as done and to retrieve its context.
type Job interface {
	// Done launches the job and returns an error if any issues occurred.
	Done() error

	// Context associated with the job.
	Context() context.Context
}

type jobImpl struct {
	handler func() error
	ctx     context.Context
}

var _ Job = (*jobImpl)(nil)

func NewJob(ctx context.Context, handler func() error) Job {
	return &jobImpl{handler, ctx}
}

func (ji *jobImpl) Done() error {
	return ji.handler()
}

func (ji *jobImpl) Context() context.Context {
	if ji.ctx == nil {
		ji.ctx = context.Background()
	}

	return ji.ctx
}
