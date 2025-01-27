package main

import (
	"context"
	"fmt"

	"github.com/nicolito128/schedule"
)

func main() {
	ctx := context.TODO()

	// Create a new PriorityQueue
	pq := schedule.NewPriorityQueue()

	// Add jobs with different priorities
	pq.Add(schedule.NewPriorityJob(func(ctx context.Context) error {
		fmt.Println("Job 1 with priority 1")
		return nil
	}, 1))

	pq.Add(schedule.NewPriorityJob(func(ctx context.Context) error {
		fmt.Println("Job 2 with priority 3")
		return nil
	}, 3))

	pq.Add(schedule.NewPriorityJob(func(ctx context.Context) error {
		fmt.Println("Job 3 with priority 2")
		return nil
	}, 2))

	// Start the priority queue execution
	pq.Run(ctx)

	// Wait for all jobs to finish
	<-pq.Done()

	// Handle errors if any
	select {
	case err := <-pq.Errors():
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	default:
		fmt.Println("All jobs completed successfully")
	}
}
