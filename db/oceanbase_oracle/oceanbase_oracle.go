package oceanbase_oracle

import (
	"database/sql"

	oceanbase "github.com/godoes/gorm-oracle"
	_ "github.com/mattn/go-oci8"
	"github.com/sunreaver/antman/v4/db"
	"gorm.io/gorm"
)

func MakeOceanbaseClient(c db.Config, gormConfig *gorm.Config) (*db.Databases, error) {
	return db.MakeClient(oceanbaseDT, c.MasterURI, c.MaxIdleConns, c.MaxOpenConns, c.LogMode, gormConfig, c.SlaveURIs...)
}

func oceanbaseDT(dsn string) gorm.Dialector {
	db, _ := sql.Open("oci8", dsn)
	return oceanbase.New(oceanbase.Config{Conn: db})
}
