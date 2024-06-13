package db

import (
	dm8 "github.com/godoes/gorm-dameng"
	oracle "github.com/godoes/gorm-oracle"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type dialector func(dsn string) gorm.Dialector

// 创建DB.
func MakeDB(c Config, gormConfig *gorm.Config) (db *Databases, err error) {
	var dt dialector
	if c.Type == DBTypeMYSQL {
		dt = mysql.Open
	} else if c.Type == DBTypeSQLite {
		dt = sqlite.Open
	} else if c.Type == DBTypePGSQL {
		dt = postgres.Open
	} else if c.Type == DBTypeORACLE {
		dt = oracle.Open
	} else if c.Type == DBTypeMogDb {
		dt = MogDBOpen
	} else if c.Type == DBTypeDameng { // 达梦
		dt = dm8.Open
	} else {
		return nil, errors.Errorf("no support db type: %v", c.Type)
	}
	return makeClient(dt, c.MasterURI, c.MaxIdleConns, c.MaxOpenConns, c.LogMode, gormConfig, c.SlaveURIs...)
}
