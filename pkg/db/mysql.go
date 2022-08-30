package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/767829413/normal-frame/internal/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// Options defines optsions for mysql database.
type Options struct {
	Host                  string
	Port                  int
	Username              string
	Password              string
	Database              string
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifeTime int
	LogLevel              int
	IsDebug               bool
}

// New create a new gorm db instance with the given options.
func New(opts *Options) (*gorm.DB, error) {
	dsn := fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=%t&loc=%s`,
		opts.Username,
		opts.Password,
		opts.Host,
		opts.Port,
		opts.Database,
		true,
		"Local")

	config := &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名，启用该选项后，`User` 表将是`user`
		},
	}

	if opts.IsDebug {
		newLogger := glogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			glogger.Config{
				SlowThreshold: time.Second,  // 慢 SQL 阈值
				LogLevel:      glogger.Info, // Log level
				Colorful:      true,         // 禁用彩色打印
			},
		)
		config.Logger = newLogger
	}

	db, err := gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		logger.LogErrorf(nil, logger.LogNameMysql, "mysql gorm.Open,error: %v", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.LogErrorf(nil, logger.LogNameMysql, "mysql db.DB(),error: %v", err)
		return nil, err
	}

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Duration(opts.MaxConnectionLifeTime) * time.Second)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(opts.MaxIdleConnections)

	return db, nil
}
