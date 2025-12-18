package common

import (
	dm8 "github.com/godoes/gorm-dameng"
	oracle "github.com/godoes/gorm-oracle"
	"github.com/pkg/errors"
	"github.com/sunreaver/antman/v4/db"
	"github.com/sunreaver/antman/v4/db/mogdb"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 创建DB.
func MakeDB(c db.Config, gormConfig *gorm.Config) (*db.Databases, error) {
	var dt db.Dialector
	if c.Type == db.DBTypeMYSQL {
		dt = mysql.Open
	} else if c.Type == db.DBTypeSQLite {
		dt = sqlite.Open
	} else if c.Type == db.DBTypePGSQL {
		dt = postgres.Open
	} else if c.Type == db.DBTypeORACLE {
		dt = oracle.Open
	} else if c.Type == db.DBTypeMogDb {
		dt = mogdb.MogDBOpen
	} else if c.Type == db.DBTypeDameng { // 达梦
		dt = dm8.Open
	} else if c.Type == db.DBTypeOBOracle {
		dt = oracle.Open
	} else {
		return nil, errors.Errorf("no support db type: %v", c.Type)
	}
	return db.MakeClient(dt, c.MasterURI, c.MaxIdleConns, c.MaxOpenConns, c.LogMode, gormConfig, c.SlaveURIs...)
}
