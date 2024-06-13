package sqlite

import (
	"github.com/sunreaver/antman/v4/db"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func MakeSqliteClient(c db.Config, gormConfig *gorm.Config) (*db.Databases, error) {
	return db.MakeClient(sqlite.Open, c.MasterURI, c.MaxIdleConns, c.MaxOpenConns, c.LogMode, gormConfig, c.SlaveURIs...)
}
