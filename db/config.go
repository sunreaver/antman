package db

// Config Config
type Config struct {
	Type         string `toml:"type"`
	URI          string `toml:"uri"`
	LogMode      bool   `toml:"log_mode"`
	MaxIdleConns int    `toml:"max_idle_conns"`
	MaxOpenConns int    `toml:"max_open_conns"`
}
