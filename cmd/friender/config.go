package main

var Config = Cfg{}

type Cfg struct {
	Config struct {
		ClearDups bool `toml:"clear_dups"`
		Threads   int  `toml:"threads"`
		MaxDm     int  `toml:"max_dm"`
	} `toml:"config"`
	Discord struct {
		Build int `toml:"build"`
	} `toml:"discord"`
	Bridge struct {
		Enable  bool   `toml:"enable"`
		Port    int    `toml:"port"`
	} `toml:"bridge"`
}
