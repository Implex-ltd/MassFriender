package instance

import u "github.com/Implex-ltd/ucdiscord/ucdiscord"

var (
	// Friend request sent
	STATUS_PROCESSED = 0

	// Can't find username (username may be invalid)
	STATUS_UNPROCESSABLE = 1

	// Token got a captcha / unknown error
	STATUS_NIL = 2

	// Token is ratelimited (http 429)
	STATUS_RATELIMIT = 3
)

type Config struct {
	MaxTask      int
	EnableSolver bool
	Client       *u.Client
}

type Report struct {
	Success, Error                    int
	Captcha, Ratelimited, InvalidUser bool
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
