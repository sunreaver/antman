package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func MakeMysqlClient(c Config, gormConfig *gorm.Config) (db *Databases, err error) {
	return makeClient(mysql.Open, c.MasterURI, c.MaxIdleConns, c.MaxOpenConns, c.LogMode, gormConfig, c.SlaveURIs...)
}
