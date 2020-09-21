package db

// Config Config.
type Config struct {
	Type         string   `toml:"type"`
	MasterURI    string   `toml:"master_uri"`
	SlaveURIs    []string `toml:"slave_uris"`
	LogMode      bool     `toml:"log_mode"`
	MaxIdleConns int      `toml:"max_idle_conns"`
	MaxOpenConns int      `toml:"max_open_conns"`
}
