package db

import (
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

type Dialector func(dsn string) gorm.Dialector

func MakeClient(dt Dialector, master string, maxIdle, maxOpen int, logMode bool, gormConfig *gorm.Config, slaves ...string) (db *Databases, err error) {
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
