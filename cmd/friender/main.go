package main

import (
	"log"
	"os"
	"time"

	"github.com/Implex-ltd/fingerprint-client/fpclient"
	"github.com/Implex-ltd/friender/internal/instance"
	"github.com/zenthangplus/goccm"
)

var (
	fp       *fpclient.Fingerprint
	PoolSize = 10
)

func GatherTasklist(length int) []string {
	out := make([]string, 0)
	
	for i := 0; i < length; i++ {
		username, err := Inputs["username"].Next()
		if err != nil {
			break
		}

		Inputs["username"].Lock(username)
		out = append(out, username)
	}

	return out
}

func ThreadWorker(token string) error {
	I, err := instance.NewInstance(&instance.Config{
		Token:   token,
		MaxTask: 5,
	})
	if err != nil {
		return err
	}

	i := 0
	for !I.IsCompleted() {
		defer func() {
			i++
		}()

		Tasklist := GatherTasklist(I.Config.MaxTask - I.Cache.Report.Success)
		if len(Tasklist) == 0 {
			break
		}

		I.PushTask(Tasklist)

		output, err := I.DoTask()
		if err != nil {
			log.Println(err)
		}

		log.Printf("[%s] job #%d done: %v, output: %v", token[:25], i, I.Cache.Report, output)

		for _, task := range output.Unprocessed {
			Inputs["username"].Unlock(task)
		}

		for _, task := range output.Unprocessable {
			Inputs["username"].Remove(task)
		}

		for _, task := range output.Processed {
			Inputs["username"].Remove(task)
		}

		Processed += len(output.Processed)
		Unprocessable += len(output.Unprocessable)
		Unprocessed = len(Inputs["username"].List)

		I.Cache.Taskout = instance.Taskout{}
	}

	log.Printf("[%s] job done: %v", token[:25], I.Cache.Report)
	return nil
}

func main() {
	if err := LoadDataset(); err != nil {
		panic(err)
	}

	var err error
	if fp, err = fpclient.LoadFingerprint(&fpclient.LoadingConfig{
		FilePath: "../../assets/chrome.json",
	}); err != nil {
		panic(err)
	}

	go ConsoleWorker()

	c := goccm.New(PoolSize)
	tokenlength := len(Inputs["tokens"].List)

	for i := 0; i < tokenlength; i++ {
		c.Wait()

		if len(Inputs["username"].List) == 0 {
			c.Done()
			continue
		}

		token, err := Inputs["tokens"].Next()
		if err != nil {
			log.Println("All tokens processed")
			break
		}

		Inputs["tokens"].Lock(token)

		go func(token string) {
			defer c.Done()

			if err := ThreadWorker(token); err != nil {
				log.Printf("[%s] %s", token[:25], err.Error())
			}
			time.Sleep(time.Millisecond * 20)
		}(token)
	}

	c.WaitAllDone()
	log.Printf("All threads exited (Unprocessed: %d, Processed: %d, Unprocessable: %d)", Unprocessed, Processed, Unprocessable)
	os.Exit(0)
}
