package instance

import u "github.com/Implex-ltd/ucdiscord/ucdiscord"

var (
	STATUS_PROCESSED     = 0
	STATUS_UNPROCESSABLE = 1
	STATUS_NIL           = 2
)

type Config struct {
	MaxTask int
	Client  *u.Client
}

type Report struct {
	Success, Error       int
	Captcha, Ratelimited bool
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
