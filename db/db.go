package db

import (
	"github.com/jinzhu/gorm"
)

// MakeDB MakeDB
func MakeDB(c Config) (*gorm.DB, error) {
	tmp, e := gorm.Open(c.Type, c.URI)
	if e != nil {
		return nil, e
	}
	tmp.LogMode(c.LogMode)
	tmp.DB().SetMaxIdleConns(c.MaxIdleConns)
	tmp.DB().SetMaxOpenConns(c.MaxOpenConns)

	return tmp, nil
}
