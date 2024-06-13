package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MakePostgresqlClient(c Config, gormConfig *gorm.Config) (db *Databases, err error) {
	return makeClient(postgres.Open, c.MasterURI, c.MaxIdleConns, c.MaxOpenConns, c.LogMode, gormConfig, c.SlaveURIs...)
}
