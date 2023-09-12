package instance

import (
	"log"
	"math/rand"
	"time"
)

func NewInstance(config *Config) (*Instance, error) {
	return &Instance{
		Config: config,
		Cache: &Cache{
			Report: Report{
				Success: 0,
				Error:   0,
				Captcha: false,
			},
		},
	}, nil
}

func (I *Instance) IsCompleted() bool {
	return I.Cache.Report.Captcha || I.Cache.Report.Success >= I.Config.MaxTask
}

func (I *Instance) PushTask(usernames []string) {
	I.Cache.Tasklist = usernames
}

func (I *Instance) Do(task string) int {
	rand.Seed(time.Now().UnixNano())
	randomValue := rand.Intn(4)

	return randomValue
}

func (I *Instance) DoTask() (*Taskout, error) {
	for i, task := range I.Cache.Tasklist {
		if I.IsCompleted() {
			break
		}

		switch I.Do(task) {
		case STATUS_PROCESSED:
			I.Cache.Report.Success++
			I.Cache.Taskout.Processed = append(I.Cache.Taskout.Processed, task)
		case STATUS_UNPROCESSABLE:
			I.Cache.Taskout.Unprocessable = append(I.Cache.Taskout.Unprocessable, task)
		default:
			I.Cache.Report.Error++
			continue
		}

		log.Printf("[#%d] doing -> %s", i, task)
	}

	for _, task := range I.Cache.Tasklist {
		if !containsTask(I.Cache.Taskout.Processed, task) && !containsTask(I.Cache.Taskout.Unprocessable, task) {
			I.Cache.Taskout.Unprocessed = append(I.Cache.Taskout.Unprocessed, task)
		}
	}

	return &I.Cache.Taskout, nil
}

// utils

func containsTask(slice []string, task string) bool {
	for _, t := range slice {
		if t == task {
			return true
		}
	}
	return false
}
