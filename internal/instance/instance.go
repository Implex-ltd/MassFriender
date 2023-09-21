package instance

import (
	"log"
	"strings"
	"sync"

	u "github.com/Implex-ltd/ucdiscord/ucdiscord"
)

func NewInstance(config *Config) (*Instance, error) {
	if _, err := config.Client.WsConnect(); err != nil {
		return nil, err
	}

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
	return I.Cache.Report.Ratelimited || I.Cache.Report.Captcha || I.Cache.Report.Success >= I.Config.MaxTask
}

func (I *Instance) PushTask(usernames []string) {
	I.Cache.Tasklist = usernames
}

func (I *Instance) Do(task string) int {
	success, captcha, err := I.Config.Client.SendFriend(&u.FriendConfig{
		Username: task,
	})

	if err != nil {
		if strings.Contains(err.Error(), "429") {
			I.Cache.Report.Ratelimited = true
			return STATUS_RATELIMIT
		}

		if strings.Contains(err.Error(), "user not found") {
			I.Cache.Report.InvalidUser = true
			return STATUS_UNPROCESSABLE
		}

		return STATUS_NIL
	}

	if captcha != nil {
		I.Cache.Report.Captcha = true
		return STATUS_NIL
	}

	log.Printf("[#%d] %s -> `%s` %v", I.Cache.Report.Success, I.Config.Client.Config.Token, task, success)
	if success {
		return STATUS_PROCESSED
	}

	return STATUS_NIL
}

func (I *Instance) DoTask() (*Taskout, error) {
	var wg sync.WaitGroup

	for _, task := range I.Cache.Tasklist {
		wg.Add(1)

		go func(task string) {
			defer wg.Done()

			switch I.Do(task) {
			case STATUS_PROCESSED:
				I.Cache.Report.Success++
				I.Cache.Taskout.Processed = append(I.Cache.Taskout.Processed, task)
			case STATUS_UNPROCESSABLE:
				I.Cache.Taskout.Unprocessable = append(I.Cache.Taskout.Unprocessable, task)
			case STATUS_RATELIMIT:
				I.Cache.Taskout.Unprocessed = append(I.Cache.Taskout.Unprocessed, task)
			default:
				I.Cache.Report.Error++
			}
		}(task)
	}

	wg.Wait()

	for _, task := range I.Cache.Tasklist {
		if !containsTask(I.Cache.Taskout.Processed, task) && !containsTask(I.Cache.Taskout.Unprocessable, task) {
			I.Cache.Taskout.Unprocessed = append(I.Cache.Taskout.Unprocessed, task)
		}
	}

	return &I.Cache.Taskout, nil
}
