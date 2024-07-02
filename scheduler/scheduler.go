package scheduler

import (
	"log"
	"time"
)

// StartScheduler starts a scheduler to run the given function every 5 minutes.
func StartScheduler(task func()) {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			log.Println("Executing scheduled task")
			task()
		}
	}()
	select {} // Block forever
}
