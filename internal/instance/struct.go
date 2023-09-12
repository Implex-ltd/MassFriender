package instance

var (
	STATUS_PROCESSED     = 0
	STATUS_UNPROCESSABLE = 1
)

type Config struct {
	MaxTask int
	Token   string
}

type Report struct {
	Success, Error int
	Captcha        bool
}

type Taskout struct {
	Processed     []string
	Unprocessed   []string
	Unprocessable []string
}

type Cache struct {
	Tasklist []string
	Taskout  Taskout
	Report   Report
}

type Instance struct {
	Config *Config
	Cache  *Cache
}
