package redis

import (
	"github.com/sunreaver/tomlanalysis/timesize"
)

// Config Config
type Config struct {
	Hosts        []string          `toml:"hosts"`
	Password     string            `toml:"password"`
	Prefix       string            `toml:"prefix"`
	DialTimeout  timesize.Duration `toml:"dial_timeout"`
	WriteTimeout timesize.Duration `toml:"write_timeout"`
	ReadTimeout  timesize.Duration `toml:"read_timeout"`
	Poolsize     int               `toml:"poolsize"`
	PoolTimeout  timesize.Duration `toml:"pool_timeout"`
	DB           int               `toml:"db"`
}
