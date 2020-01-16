package db

import (
	"sync/atomic"

	"github.com/jinzhu/gorm"
)

type Databases struct {
	master     *gorm.DB
	slaves     []*gorm.DB
	slaveCount uint32 // 从 db 个数
	slaveIndex uint32 // 从 db 当前 index
}

func (db *Databases) Free() {
	if db.master != nil {
		db.master.Close()
	}
	slaves := db.slaves
	db.slaveCount = 0
	db.slaveIndex = 0
	db.slaves = nil
	for _, s := range slaves {
		s.Close()
	}
}

func (db *Databases) Master() *gorm.DB {
	return db.master
}

func (db *Databases) Slave() *gorm.DB {
	if db.slaveCount == 0 {
		return db.Master()
	}
	index := atomic.AddUint32(&db.slaveIndex, 1)
	index %= db.slaveCount
	return db.slaves[index]
}
