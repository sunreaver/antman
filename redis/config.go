package redis

// Config Config
type Config struct {
	Hosts        []string `toml:"hosts"`
	Prefix       string   `toml:"prefix"`
	DialTimeout  int      `toml:"dial_timeout"`
	WriteTimeout int      `toml:"write_timeout"`
	ReadTimeout  int      `toml:"read_timeout"`
	Poolsize     int      `toml:"poolsize"`
	PoolTimeout  int      `toml:"pool_timeout"`
}
