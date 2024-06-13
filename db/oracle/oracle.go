package oracle

import (
	oracle "github.com/godoes/gorm-oracle"
	"github.com/sunreaver/antman/v4/db"
	"gorm.io/gorm"
)

func MakeOracleClient(c db.Config, gormConfig *gorm.Config) (*db.Databases, error) {
	return db.MakeClient(oracle.Open, c.MasterURI, c.MaxIdleConns, c.MaxOpenConns, c.LogMode, gormConfig, c.SlaveURIs...)
}
