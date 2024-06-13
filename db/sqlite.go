package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func MakeSqliteClient(c Config, gormConfig *gorm.Config) (db *Databases, err error) {
	return makeClient(sqlite.Open, c.MasterURI, c.MaxIdleConns, c.MaxOpenConns, c.LogMode, gormConfig, c.SlaveURIs...)
}
