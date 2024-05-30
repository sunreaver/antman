package db

import (
	"time"

	dm8 "github.com/godoes/gorm-dameng"
	oracle "github.com/godoes/gorm-oracle"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
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

func MakeMysqlClient(c Config, gormConfig *gorm.Config) (db *Databases, err error) {
	return makeClient(mysql.Open, c.MasterURI, c.MaxIdleConns, c.MaxOpenConns, c.LogMode, gormConfig, c.SlaveURIs...)
}

func MakeSqliteClient(c Config, gormConfig *gorm.Config) (db *Databases, err error) {
	return makeClient(sqlite.Open, c.MasterURI, c.MaxIdleConns, c.MaxOpenConns, c.LogMode, gormConfig, c.SlaveURIs...)
}

func MakeOracleClient(c Config, gormConfig *gorm.Config) (db *Databases, err error) {
	return makeClient(oracle.Open, c.MasterURI, c.MaxIdleConns, c.MaxOpenConns, c.LogMode, gormConfig, c.SlaveURIs...)
}

func MakePostgresqlClient(c Config, gormConfig *gorm.Config) (db *Databases, err error) {
	return makeClient(postgres.Open, c.MasterURI, c.MaxIdleConns, c.MaxOpenConns, c.LogMode, gormConfig, c.SlaveURIs...)
}

func MakeMogdbClient(c Config, gormConfig *gorm.Config) (db *Databases, err error) {
	return makeClient(MogDBOpen, c.MasterURI, c.MaxIdleConns, c.MaxOpenConns, c.LogMode, gormConfig, c.SlaveURIs...)
}

func MakeDamengClient(c Config, gormConfig *gorm.Config) (db *Databases, err error) {
	return makeClient(dm8.Open, c.MasterURI, c.MaxIdleConns, c.MaxOpenConns, c.LogMode, gormConfig, c.SlaveURIs...)
}

func makeClient(dt dialector, master string, maxIdle, maxOpen int, logMode bool, gormConfig *gorm.Config, slaves ...string) (db *Databases, err error) {
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
	}).SetConnMaxIdleTime(time.Second * 5).
		SetConnMaxLifetime(time.Second * 20).
		SetMaxIdleConns(maxIdle).
		SetMaxOpenConns(maxOpen))

	if e != nil {
		return nil, errors.Wrap(e, "set dbresolver")
	}

	db = &Databases{
		master: tmp,
	}

	return db, nil
}
