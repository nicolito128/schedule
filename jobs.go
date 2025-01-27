package schedule

import "context"

// Job represents a unit of work that can be executed.
type Job interface {
	// Done launches the job and returns an error if any issues occurred.
	Done(ctx context.Context) error
}

type jobImpl struct {
	handler func(context.Context) error
}

var _ Job = (*jobImpl)(nil)

func NewJob(handler func(context.Context) error) Job {
	return &jobImpl{handler}
}

func (ji *jobImpl) Done(ctx context.Context) error {
	return ji.handler(ctx)
}
