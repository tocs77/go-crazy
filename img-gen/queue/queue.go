package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-queue/queue"
	"github.com/golang-queue/queue/core"
)

type jobData struct {
	Name    string
	Message string
}

func (j *jobData) Bytes() []byte {
	fmt.Printf("%s:%s\n", j.Name, j.Message)
	res := sleepSomeTime()
	j.Name = "I am awake"
	j.Message = res
	b, _ := json.Marshal(j)
	return b
}

func (j *jobData) Payload() []byte {
	return j.Bytes()
}

func sleepSomeTime() string {
	seconds := rand.Intn(20)
	sleepTime := time.Duration(seconds) * time.Second
	time.Sleep(sleepTime)
	return fmt.Sprintf("Commander, I slept: %d seconds", seconds)
}

func StartQueue() {
	taskN := 100
	rets := make(chan string, taskN)
	q := queue.NewPool(30, queue.WithFn(func(ctx context.Context, m core.TaskMessage) error {
		var v jobData
		if err := json.Unmarshal(m.Payload(), &v); err != nil {
			return err
		}

		// Check if v is empty after unmarshalling
		if v.Name == "" && v.Message == "" {
			fmt.Println("Warning: jobData is empty after unmarshalling")
			rets <- "Empty job data received"
		} else {
			rets <- "Hello, " + v.Name + ", " + v.Message
		}
		return nil
	}))
	defer q.Release()
	for i := 0; i < taskN; i++ {
		go func(i int) {
			q.Queue(&jobData{
				Name:    "Sleeping Gophers",
				Message: fmt.Sprintf("Hello commander, I am handling the job: %d", +i),
			})
		}(i)
	}
	for i := 0; i < taskN; i++ {
		fmt.Println("message:", <-rets)
		time.Sleep(10 * time.Millisecond)
	}
}
