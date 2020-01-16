package db

import (
	"github.com/jinzhu/gorm"
)

// MakeDB MakeDB
func MakeDB(c Config) (db *Databases, err error) {
	tmp, e := gorm.Open(c.Type, c.MasterURI)
	if e != nil {
		return nil, e
	}
	tmp.LogMode(c.LogMode)
	tmp.DB().SetMaxIdleConns(c.MaxIdleConns)
	tmp.DB().SetMaxOpenConns(c.MaxOpenConns)

	db = &Databases{
		master: tmp,
		slaves: []*gorm.DB{},
	}

	defer func() {
		if err != nil {
			db.Free()
		}
	}()

	for _, uri := range c.SlaveURIs {
		tmp, e := gorm.Open(c.Type, uri)
		if e != nil {
			return nil, e
		}
		tmp.LogMode(c.LogMode)
		tmp.DB().SetMaxIdleConns(c.MaxIdleConns)
		tmp.DB().SetMaxOpenConns(c.MaxOpenConns)
		db.slaves = append(db.slaves, tmp)
		db.slaveCount++
	}

	return db, nil
}
