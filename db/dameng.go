package db

import (
	dm8 "github.com/godoes/gorm-dameng"
	"gorm.io/gorm"
)

func MakeDamengClient(c Config, gormConfig *gorm.Config) (db *Databases, err error) {
	return makeClient(dm8.Open, c.MasterURI, c.MaxIdleConns, c.MaxOpenConns, c.LogMode, gormConfig, c.SlaveURIs...)
}
