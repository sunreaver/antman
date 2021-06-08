package db

import (
	"time"

	"github.com/cengsin/oracle"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

// 创建DB.
func MakeDB(c Config, gormConfig *gorm.Config) (db *Databases, err error) {
	dbMaster, e := makeClient(c.Type, c.MasterURI, c.MaxIdleConns, c.MaxOpenConns, c.LogMode, gormConfig, c.SlaveURIs...)
	if e != nil {
		return nil, errors.Wrap(e, "make db")
	}

	db = &Databases{
		master: dbMaster,
	}

	return db, nil
}

type dialector func(dsn string) gorm.Dialector

func makeClient(dbType string, master string, maxIdle, maxOpen int, logMode bool, gormConfig *gorm.Config, slaves ...string) (*gorm.DB, error) {
	var dt dialector
	if dbType == "mysql" {
		dt = mysql.Open
	} else if dbType == "sqlite" {
		dt = sqlite.Open
	} else if dbType == "postgres" {
		dt = postgres.Open
	} else if dbType == "oracle" {
		dt = oracle.Open
	} else {
		return nil, errors.Errorf("no support db type: %v", dbType)
	}
	var loggerMode logger.Interface
	if logMode {
		loggerMode = logger.Default.LogMode(logger.Info)
	} else {
		loggerMode = logger.Default.LogMode(logger.Silent)
	}
	if gormConfig == nil {
		gormConfig = &gorm.Config{
			Logger: loggerMode,
		}
	} else if gormConfig.Logger == nil {
		gormConfig.Logger = loggerMode
	}
	tmp, e := gorm.Open(dt(master), gormConfig)
	if e != nil {
		return nil, errors.Wrap(e, "gorm open")
	}

	// 组成slave
	ss := make([]gorm.Dialector, len(slaves))
	for index, s := range slaves {
		ss[index] = dt(s)
	}

	e = tmp.Use(dbresolver.Register(dbresolver.Config{
		Replicas: ss,
		Policy:   dbresolver.RandomPolicy{},
	}).SetConnMaxIdleTime(time.Minute * 5).
		SetConnMaxLifetime(time.Hour * 24).
		SetMaxIdleConns(maxIdle).
		SetMaxOpenConns(maxOpen))

	if e != nil {
		return nil, errors.Wrap(e, "set dbresolver")
	}

	return tmp, nil
}
