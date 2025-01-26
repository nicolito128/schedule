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

	lifo := schedule.NewLIFO()

	// Adding jobs to the LIFO scheduler
	for i := 1; i <= 5; i++ {
		jobID := i
		lifo.Add(schedule.NewJob(ctx, func() error {
			fmt.Printf("Executing job %d\n", jobID)
			time.Sleep(1 * time.Second)
			return nil
		}))
	}

	// Running the LIFO scheduler
	lifo.Run(ctx)

	// Waiting for all jobs to complete
	<-lifo.Done()

	// Checking for errors
	select {
	case err := <-lifo.Errors():
		if err != nil {
			fmt.Printf("Error encountered: %v\n", err)
		}
	default:
		fmt.Println("All jobs completed successfully.")
	}
}
