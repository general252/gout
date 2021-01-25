package utask

import (
	"log"
	"time"
)

func ExampleAddTaskEvery() {
	defer StopCron()
	_, _ = AddTaskEvery(time.Second, func() {
		log.Println("1 second ", time.Now())
	})

	_, _ = AddTaskEvery(time.Second*5, func() {
		log.Println("5 second ", time.Now())
	})

	time.Sleep(time.Minute)

	// output:
	//
}
