package store

import (
	"fmt"
	"sync"

	mylog "github.com/767829413/normal-frame/internal/pkg/logger"
	"github.com/767829413/normal-frame/internal/pkg/options"
	"github.com/767829413/normal-frame/pkg/db"
	"gorm.io/gorm"
)

var (
	dbHandler *datastore
	once      sync.Once
)

type datastore struct {
	db *gorm.DB

	// can include two database instance if needed
	// docker *grom.DB
	// db *gorm.DB
}

// GetMySQLIncOr create mysql factory with the given config.
func GetMySQLIncOr(opts *options.MySQLOptions) *datastore {
	if dbHandler != nil {
		return dbHandler
	}
	if opts != nil && !opts.Enabled {
		return nil
	}
	if opts == nil && dbHandler == nil {
		return nil
	}
	var err error
	var dbIns *gorm.DB
	once.Do(func() {
		options := &db.Options{
			Host:                  opts.Host,
			Username:              opts.Username,
			Password:              opts.Password,
			Database:              opts.Database,
			MaxIdleConnections:    opts.MaxIdleConnections,
			MaxOpenConnections:    opts.MaxOpenConnections,
			MaxConnectionLifeTime: opts.MaxConnectionLifeTime,
			LogLevel:              opts.LogLevel,
		}
		dbIns, err = db.New(options)
		dbHandler = &datastore{dbIns}
	})
	if err != nil {
		panic(fmt.Sprintf("GetMySQLIncOr err : %v", err))
	}
	return dbHandler

}

func (d *datastore) GetDb() *gorm.DB {
	return d.db
}

func (d *datastore) Close() error {
	if dbHandler.db != nil {
		db, err := dbHandler.db.DB()
		if err != nil {
			mylog.LogError(nil, mylog.LogNameMysql, "dbHandler get gorm db instance failed")
		}
		return db.Close()
	}
	return nil
}
