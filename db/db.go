package db

import (
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 创建DB.
func MakeDB(c Config, gormConfig *gorm.Config) (db *Databases, err error) {
	dbMaster, e := makeClient(c.Type, c.MasterURI, c.MaxIdleConns, c.MaxOpenConns, c.LogMode, gormConfig)
	if e != nil {
		return nil, errors.Wrap(e, "master")
	}

	db = &Databases{
		master: dbMaster,
		slaves: []*gorm.DB{},
	}

	defer func() {
		if err != nil {
			db.Free()
		}
	}()

	for _, uri := range c.SlaveURIs {
		dbSlave, e := makeClient(c.Type, uri, c.MaxIdleConns, c.MaxOpenConns, c.LogMode, gormConfig)
		if e != nil {
			return nil, e
		}
		db.slaves = append(db.slaves, dbSlave)
		db.slaveCount++
	}

	return db, nil
}

func makeClient(dbType string, uri string, maxIdle, maxOpen int, logMode bool, gormConfig *gorm.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector
	if dbType == "mysql" {
		dialector = mysql.Open(uri)
	} else if dbType == "sqlite" {
		dialector = sqlite.Open(uri)
	} else if dbType == "postgres" {
		dialector = postgres.Open(uri)
	} else {
		return nil, errors.Errorf("no support db type: %v", dbType)
	}
	var loggerMode logger.Interface
	if logMode {
		loggerMode = logger.Default.LogMode(logger.Silent)
	} else {
		loggerMode = logger.Default
	}
	if gormConfig == nil {
		gormConfig = &gorm.Config{
			Logger: loggerMode,
		}
	} else if gormConfig.Logger == nil {
		gormConfig.Logger = loggerMode
	}
	tmp, e := gorm.Open(dialector, gormConfig)
	if e != nil {
		return nil, e
	}
	dbTmp, e := tmp.DB()
	if e != nil {
		return nil, errors.Wrap(e, "master")
	}
	dbTmp.SetMaxIdleConns(maxIdle)
	dbTmp.SetMaxOpenConns(maxOpen)

	return tmp, nil
}
