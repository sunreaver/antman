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

func (md *MogDBDialector) Initialize(db *gorm.DB) error {
	if err := md.Dialector.Initialize(db); err != nil {
		return err
	}
	// callbacks.RegisterDefaultCallbacks(db, &callbacks.Config{
	// 	WithReturning: false,
	// })
	return nil
}
