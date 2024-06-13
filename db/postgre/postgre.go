package postgre

import (
	"github.com/sunreaver/antman/v4/db"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MakePostgresqlClient(c db.Config, gormConfig *gorm.Config) (*db.Databases, error) {
	return db.MakeClient(postgres.Open, c.MasterURI, c.MaxIdleConns, c.MaxOpenConns, c.LogMode, gormConfig, c.SlaveURIs...)
}
