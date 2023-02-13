package main

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

const TasksCount = 15 // you can edit this

type Task struct {
	id          uint32
	createdAt   time.Time
	completedAt time.Time
	success     bool
	error       error
}

func (t Task) IsCorrect() bool {
	return !t.createdAt.IsZero()
}

func (t Task) Run() Task {
	t.success = t.IsCorrect()
	if !t.success {
		t.error = fmt.Errorf("the moon is in the wrong phase")
	}

	time.Sleep(time.Millisecond * 150)

	t.completedAt = time.Now()

	return t
}

func seedTasks(c chan<- Task, count int) {
	for i := 0; i < count; i++ {
		go func() {
			c <- generateTask()
		}()
	}
}

func generateTask() Task {
	now := time.Now()
	id := uuid.New().ID()

	createdAt := now
	if shouldAddIncorrectTask(id) {
		createdAt = time.Time{}
	}

	time.Sleep(time.Millisecond * 10) // to avoid duplicate ids

	return Task{createdAt: createdAt, id: id}
}

func shouldAddIncorrectTask(id uint32) bool {
	// Original code was t.Nanosecond()%2 > 0, but chance was very small

	return id%2 > 0
}

func main() {
	taskQueue := make(chan Task, 10)

	go seedTasks(taskQueue, TasksCount)

	successfulTasks := make(chan Task)
	failedTasks := make(chan Task)

	sortCompletedTask := func(t Task) {
		if t.success {
			successfulTasks <- t
		} else {
			failedTasks <- t
		}
	}

	go func() {
		for t := range taskQueue {
			go func(t Task) {
				t = t.Run()

				sortCompletedTask(t)
			}(t)
		}

		close(taskQueue)
	}()

	//var wg sync.WaitGroup
	//wg.Add(2)

	go printSuccessfulTasks(successfulTasks /*, &wg*/)
	go printFailedTasks(failedTasks /*, &wg*/)

	//wg.Wait()

	time.Sleep(time.Second * 1) // I'm going to replace this with wait groups
}

func printSuccessfulTasks(tasks chan Task /*, wg *sync.WaitGroup*/) {
	//defer wg.Done()

	for t := range tasks {
		fmt.Printf("task #%d completed at %s\n", t.id, t.completedAt.Format(time.RFC3339))
	}

	close(tasks)
}

func printFailedTasks(tasks chan Task /*, wg *sync.WaitGroup*/) {
	//defer wg.Done()

	for t := range tasks {
		fmt.Printf(
			"task #%d failed at %s with error \"%s\"\n",
			t.id,
			t.completedAt.Format(time.RFC3339),
			t.error,
		)
	}

	close(tasks)
}
