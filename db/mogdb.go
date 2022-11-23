package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MogDBDialector struct{ *postgres.Dialector }

func MogDBOpen(dsn string) gorm.Dialector {
	dial := postgres.New(postgres.Config{
		DriverName: "opengauss",
		DSN:        dsn,
	}).(*postgres.Dialector)
	return &MogDBDialector{dial}
}

func (md *MogDBDialector) Name() string {
	return "mogdb"
}
