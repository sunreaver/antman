package dameng

import (
	dm8 "github.com/godoes/gorm-dameng"
	"github.com/sunreaver/antman/v4/db"
	"gorm.io/gorm"
)

func MakeDamengClient(c db.Config, gormConfig *gorm.Config) (*db.Databases, error) {
	return db.MakeClient(dm8.Open, c.MasterURI, c.MaxIdleConns, c.MaxOpenConns, c.LogMode, gormConfig, c.SlaveURIs...)
}
