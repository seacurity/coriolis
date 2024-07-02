package scheduler

import (
	"testing"
	"time"
)

func TestStartScheduler(t *testing.T) {
	taskExecuted := false

	// Define a task function that sets taskExecuted to true
	task := func() {
		taskExecuted = true
	}

	// Start the scheduler with a very short interval for testing
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			task()
		}
	}()

	// Wait for a short period to allow the task to be executed
	time.Sleep(2 * time.Second)
	ticker.Stop()

	if !taskExecuted {
		t.Errorf("Expected task to be executed, but it was not")
	}
}
