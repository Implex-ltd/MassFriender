package console

import (
	"fmt"
	"github.com/gookit/color"
	"time"

	"github.com/Implex-ltd/friender/internal/utils"
)

var (
	Processed, Unprocessable, Unprocessed, Ratelimit, Captcha int
	TotalArr                                                  []int
)

func SetTitle(title string) {
	fmt.Printf("\033]0;%s\007", title)
}

func ConsoleWorker() {
	t := time.NewTicker(time.Millisecond)

	go func() {
		for {
			select {
			case <-t.C:
				SetTitle(fmt.Sprintf("Unprocessed: %d, Processed: %d, Unprocessable: %d, Ratelimit: %d, Captcha: %d, Avg: %.2f", Unprocessed, Processed, Unprocessable, Ratelimit, Captcha, utils.CalculateAverage(TotalArr)))
			}
		}
	}()
}

func Log(content string) {
	color.Println(fmt.Sprintf("[%s] %s", time.Now().Format("15:04:05"), content))
}
