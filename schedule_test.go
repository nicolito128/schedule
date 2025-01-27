package schedule_test

import (
	"context"
	"testing"
	"time"

	"github.com/nicolito128/schedule"
)

func TestLIFO(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	lifo := schedule.NewLIFO()

	for i := 1; i <= 5; i++ {
		jobID := i
		lifo.Add(schedule.NewJob(func(ctx context.Context) error {
			t.Logf("Executing LIFO job %d\n", jobID)
			return nil
		}))
	}

	lifo.Run(ctx)
	<-lifo.Done()

	select {
	case err := <-lifo.Errors():
		if err != nil {
			t.Errorf("Error encountered in LIFO: %v", err)
		}
	default:
		t.Log("All LIFO jobs completed successfully.")
	}
}

func TestFIFO(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fifo := schedule.NewFIFO()

	for i := 1; i <= 5; i++ {
		jobID := i
		fifo.Add(schedule.NewJob(func(ctx context.Context) error {
			t.Logf("Executing FIFO job %d\n", jobID)
			return nil
		}))
	}

	fifo.Run(ctx)
	<-fifo.Done()

	select {
	case err := <-fifo.Errors():
		if err != nil {
			t.Errorf("Error encountered in FIFO: %v", err)
		}
	default:
		t.Log("All FIFO jobs completed successfully.")
	}
}

func TestPriorityQueue(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pq := schedule.NewPriorityQueue()

	for i := 1; i <= 5; i++ {
		jobID := i
		priority := 5 - i // Higher priority for lower jobID
		pq.Add(schedule.NewPriorityJob(func(ctx context.Context) error {
			t.Logf("Executing PriorityQueue job %d with priority %d\n", jobID, priority)
			return nil
		}, priority))
	}

	pq.Run(ctx)
	<-pq.Done()

	select {
	case err := <-pq.Errors():
		if err != nil {
			t.Errorf("Error encountered in PriorityQueue: %v", err)
		}
	default:
		t.Log("All PriorityQueue jobs completed successfully.")
	}
}
