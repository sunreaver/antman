package db

import (
	"sync/atomic"

	"gorm.io/gorm"
)

type Databases struct {
	master     *gorm.DB
	slaves     []*gorm.DB
	slaveCount uint32 // 从 db 个数
	slaveIndex uint32 // 从 db 当前 index
}

func (db *Databases) Free() {
	if db.master != nil {
		master, e := db.master.DB()
		if e == nil {
			master.Close()
		}
	}
	db.master = nil
	slaves := db.slaves
	db.slaveCount = 0
	db.slaveIndex = 0
	db.slaves = nil
	for _, s := range slaves {
		slave, e := s.DB()
		if e == nil {
			slave.Close()
		}
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
