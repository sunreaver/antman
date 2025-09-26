package oracle

import (
	"database/sql"

	oracle "github.com/godoes/gorm-oracle"
	"github.com/sijms/go-ora/v2"
	"github.com/sunreaver/antman/v4/db"
	"gorm.io/gorm"
)

func MakeOracleClient(c db.Config, gormConfig *gorm.Config) (*db.Databases, error) {
	return db.MakeClient(oracle.Open, c.MasterURI, c.MaxIdleConns, c.MaxOpenConns, c.LogMode, gormConfig, c.SlaveURIs...)
}

func MakeOracleClientDriverByGoOra(c db.Config, gormConfig *gorm.Config) (*db.Databases, error) {
	return db.MakeClient(gooraOracleDT, c.MasterURI, c.MaxIdleConns, c.MaxOpenConns, c.LogMode, gormConfig, c.SlaveURIs...)
}

/*
参考：https://github.com/sijms/go-ora?tab=readme-ov-file#connect-to-multiple-database
sql.Open使用默认驱动，使用于单个项目的数据库
要使用多个数据库，应该为每个数据库创建单独的驱动程序
*/
func gooraOracleDT(dsn string) gorm.Dialector {
	connector := go_ora.NewConnector(dsn)
	db := sql.OpenDB(connector)
	return oracle.New(oracle.Config{Conn: db})
}
