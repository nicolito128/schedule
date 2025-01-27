package main

import (
	"context"
	"fmt"
	"time"

	"github.com/nicolito128/schedule"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fifo := schedule.NewFIFO()

	// Adding jobs to the FIFO scheduler
	for i := 1; i <= 5; i++ {
		jobID := i
		fifo.Add(schedule.NewJob(func(ctx context.Context) error {
			fmt.Printf("Executing job %d\n", jobID)
			time.Sleep(1 * time.Second)
			return nil
		}))
	}

	// Running the FIFO scheduler
	fifo.Run(ctx)

	// Waiting for all jobs to complete
	<-fifo.Done()

	// Checking for errors
	select {
	case err := <-fifo.Errors():
		if err != nil {
			fmt.Printf("Error encountered: %v\n", err)
		}
	default:
		fmt.Println("All jobs completed successfully.")
	}
}
