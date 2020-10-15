package db

import (
	"gorm.io/gorm"
)

type Databases struct {
	master *gorm.DB
}

func (db *Databases) Free() {
	db.master = nil
}

func (db *Databases) Master() *gorm.DB {
	return db.master
}

func (db *Databases) Slave() *gorm.DB {
	return db.master
}
