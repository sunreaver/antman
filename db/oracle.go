package db

import (
	oracle "github.com/godoes/gorm-oracle"
	"gorm.io/gorm"
)

func MakeOracleClient(c Config, gormConfig *gorm.Config) (db *Databases, err error) {
	return makeClient(oracle.Open, c.MasterURI, c.MaxIdleConns, c.MaxOpenConns, c.LogMode, gormConfig, c.SlaveURIs...)
}
