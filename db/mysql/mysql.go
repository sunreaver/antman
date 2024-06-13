package mysql

import (
	"github.com/sunreaver/antman/v4/db"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func MakeMysqlClient(c db.Config, gormConfig *gorm.Config) (*db.Databases, error) {
	return db.MakeClient(mysql.Open, c.MasterURI, c.MaxIdleConns, c.MaxOpenConns, c.LogMode, gormConfig, c.SlaveURIs...)
}
