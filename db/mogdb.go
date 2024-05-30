package db

import (
	_ "gitee.com/opengauss/openGauss-connector-go-pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MogDBDialector struct{ *postgres.Dialector }

func MogDBOpen(dsn string) gorm.Dialector {
	dial := postgres.New(postgres.Config{
		DriverName:       "opengauss",
		DSN:              dsn,
		WithoutReturning: false,
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
	// 	WithReturning: true,
	// })
	return nil
}
