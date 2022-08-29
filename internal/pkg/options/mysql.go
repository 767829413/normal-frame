package options

import (
	"time"

	"github.com/spf13/pflag"
	glogger "gorm.io/gorm/logger"
)

type MySQLOptions struct {
	Enabled               bool          `json:"enabled" mapstructure:"enabled" yaml:"enabled"`
	IsDebug               bool          `json:"is-debug" mapstructure:"is-debug" yaml:"is-debug"`
	Host                  string        `json:"host,omitempty" mapstructure:"host" yaml:"host"`
	Port                  int           `json:"port,omitempty" mapstructure:"port" yaml:"port"`
	Username              string        `json:"username,omitempty" mapstructure:"username" yaml:"username"`
	Password              string        `json:"-" mapstructure:"password" yaml:"password"`
	Database              string        `json:"database"  mapstructure:"database" yaml:"database"`
	MaxIdleConnections    int           `json:"max-idle-connections,omitempty" mapstructure:"max-idle-connections" yaml:"max-idle-connections"`
	MaxOpenConnections    int           `json:"max-open-connections,omitempty" mapstructure:"max-open-connections" yaml:"max-open-connections"`
	MaxConnectionLifeTime time.Duration `json:"max-connection-life-time,omitempty" mapstructure:"max-connection-life-time" yaml:"max-connection-life-time"`
	LogLevel              int           `json:"log-level" mapstructure:"log-level" yaml:"log-level"`
}

func NewMySQLOptions() *MySQLOptions {
	return &MySQLOptions{
		Enabled:               false,
		IsDebug:               false,
		Host:                  "127.0.0.1",
		Port:                  3306,
		Username:              "",
		Password:              "",
		Database:              "",
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifeTime: time.Duration(10) * time.Second,
		LogLevel:              int(glogger.Info),
	}
}

func (o *MySQLOptions) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&o.Enabled, "mysql.enabled", o.Enabled, "Whether to enable MySQL.")

	fs.BoolVar(&o.IsDebug, "mysql.is-debug", o.IsDebug, "Whether to enable MySQL debug mode.")

	fs.StringVar(&o.Host, "mysql.host", o.Host, ""+
		"MySQL service host address. If left blank, the following related mysql options will be ignored.")

	fs.IntVar(&o.Port, "mysql.port", o.Port, ""+
		"The port MySQL is listening on.")

	fs.StringVar(&o.Username, "mysql.username", o.Username, ""+
		"Username for access to mysql service.")

	fs.StringVar(&o.Password, "mysql.password", o.Password, ""+
		"Password for access to mysql, should be used pair with password.")

	fs.StringVar(&o.Database, "mysql.database", o.Database, ""+
		"Database name for the server to use.")

	fs.IntVar(&o.MaxIdleConnections, "mysql.max-idle-connections", o.MaxOpenConnections, ""+
		"Maximum idle connections allowed to connect to mysql.")

	fs.IntVar(&o.MaxOpenConnections, "mysql.max-open-connections", o.MaxOpenConnections, ""+
		"Maximum open connections allowed to connect to mysql.")

	fs.DurationVar(&o.MaxConnectionLifeTime, "mysql.max-connection-life-time", o.MaxConnectionLifeTime, ""+
		"Maximum connection life time allowed to connect to mysql.")

	fs.IntVar(&o.LogLevel, "mysql.log-level", o.LogLevel, ""+
		"Specify gorm log level. Silent-1 Error-2 Warn-3 Info-4")
}
