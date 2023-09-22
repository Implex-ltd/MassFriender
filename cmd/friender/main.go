package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Implex-ltd/cleanhttp/cleanhttp"
	"github.com/Implex-ltd/fingerprint-client/fpclient"
	"github.com/Implex-ltd/friender/internal/console"
	"github.com/Implex-ltd/friender/internal/instance"
	"github.com/Implex-ltd/friender/internal/utils"
	u "github.com/Implex-ltd/ucdiscord/ucdiscord"
	"github.com/zenthangplus/goccm"
)

var (
	fp         *fpclient.Fingerprint
	PoolSize   = 350
	DmPerToken = 15
)

func GatherTasklist(length int) []string {
	out := make([]string, 0)

	for i := 0; i < length; i++ {
		username, err := Inputs["username"].Next()
		if err != nil {
			break
		}

		if Inputs["done"].IsInList(username) || Inputs["blacklist"].IsInList(username) {
			i--
			continue
		}

		Inputs["username"].Lock(username)
		out = append(out, username)
	}

	return out
}

func ThreadWorker(token string) error {
	proxy, err := Inputs["proxies"].Next()
	if err != nil {
		return err
	}

	http, err := cleanhttp.NewCleanHttpClient(&cleanhttp.Config{
		BrowserFp: fp,
		Proxy:     "http://" + proxy,
	})
	if err != nil {
		return err
	}

	wss, err := u.NewWebsocket(token, &u.XProp{
		BrowserVersion:    http.BaseHeader.UaInfo.BrowserVersion,
		Browser:           http.BaseHeader.UaInfo.BrowserName,
		OsVersion:         http.BaseHeader.UaInfo.OSVersion,
		Os:                http.BaseHeader.UaInfo.OSName,
		BrowserUserAgent:  fp.Navigator.UserAgent,
		ReleaseChannel:    "stable",
		SystemLocale:      "fr-FR",
		ClientBuildNumber: 226220,
		Device:            "",
	})
	if err != nil {
		return err
	}

	client, err := u.NewClient(&u.Config{
		Token:      token,
		GetCookies: true,
		Build:      226220,
		Http:       http,
		Ws:         wss,
	})
	if err != nil {
		return err
	}

	locked, _, err := client.IsLocked()
	if err != nil {
		return err
	}

	if locked {
		go utils.AppendLineInDirectory("../../assets/data", "dead.txt", token)
		return fmt.Errorf("token is locked")
	}

	I, err := instance.NewInstance(&instance.Config{
		Client:  client,
		MaxTask: DmPerToken,
	})
	if err != nil {
		go utils.AppendLineInDirectory("../../assets/data", "dead.txt", token)
		return fmt.Errorf("token is dead (%s)", err.Error())
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

		// process task that raise error
		for _, task := range output.Unprocessed {
			Inputs["username"].Unlock(task)
		}

		// process done usernames
		for _, task := range output.Processed {
			Inputs["username"].Remove(task)
			go console.Log(fmt.Sprintf(`[<fg=58cee8>%d</>] [<fg=9572e8>%d</>] [<fg=7bcf32>%s</>] sent to: %s`, console.Processed, I.Cache.Report.Success, token[:25], task))
			go utils.AppendLineInDirectory("../../assets/data", "done.txt", task)
		}

		// process invalid usernames
		for _, task := range output.Unprocessable {
			Inputs["username"].Remove(task)
			go utils.AppendLineInDirectory("../../assets/data", "blacklist.txt", task)
		}

		console.Processed += len(output.Processed)
		console.Unprocessable += len(output.Unprocessable)
		console.Unprocessed = len(Inputs["username"].List)

		I.Cache.Taskout = instance.Taskout{}
	}

	/*if !I.Cache.Report.Captcha && !I.Cache.Report.Ratelimited {
		go utils.AppendLineInDirectory("../../assets/data", "dead.txt", token)
	}*/

	defer client.Ws.Close()

	if I.Cache.Report.Success != 0 {
		console.TotalArr = append(console.TotalArr, I.Cache.Report.Success)
	}

	if I.Cache.Report.Captcha {
		console.Log(fmt.Sprintf("[<fg=58cee8>%d</>] [<fg=9572e8>%d</>] [<fg=f23a46>%s</>] captcha", console.Processed, I.Cache.Report.Success, token[:25]))
		go utils.AppendLineInDirectory("../../assets/data", "captcha.txt", token)
		console.Captcha++
		return nil
	}

	if I.Cache.Report.Ratelimited {
		console.Log(fmt.Sprintf("[<fg=58cee8>%d</>] [<fg=9572e8>%d</>] [<fg=f2a649>%s</>] ratelimit", console.Processed, I.Cache.Report.Success, token[:25]))
		go utils.AppendLineInDirectory("../../assets/data", "ratelimit.txt", token)
		console.Ratelimit++
		return nil
	}

	console.Log(fmt.Sprintf("[<fg=58cee8>%d</>] [<fg=9572e8>%d</>] [<fg=9bebd1>%s</>] job done", console.Processed, I.Cache.Report.Success, token[:25]))
	return nil
}

func main() {
	if err := LoadDataset(); err != nil {
		panic(err)
	}

	var err error
	if fp, err = fpclient.LoadFingerprint(&fpclient.LoadingConfig{
		FilePath: "../../assets/data/chrome.json",
	}); err != nil {
		panic(err)
	}

	go console.ConsoleWorker()

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

		if Inputs["dead"].IsInList(token) || Inputs["captcha"].IsInList(token) || Inputs["ratelimit"].IsInList(token) {
			c.Done()
			continue
		}

		go func(token string) {
			defer c.Done()

			if err := ThreadWorker(token); err != nil {
				console.Log(fmt.Sprintf("[%s] %s", token[:25], err.Error()))
			}
			time.Sleep(time.Millisecond * 20)
		}(token)
	}

	c.WaitAllDone()
	console.Log(fmt.Sprintf("Unprocessed: %d, Processed: %d, Unprocessable: %d, Ratelimit: %d, Captcha: %d, Avg: %.2f", console.Unprocessed, console.Processed, console.Unprocessable, console.Ratelimit, console.Captcha, utils.CalculateAverage(console.TotalArr)))
	os.Exit(0)
}
