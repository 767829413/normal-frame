package options

import (
	"github.com/spf13/pflag"
)

type RedisOptions struct {
	Enabled bool   `json:"enabled" mapstructure:"enabled" yaml:"enabled"`
	Address string `mapstructure:"address" json:"address" yaml:"address"`
	Prefix  string `json:"prefix,omitempty" mapstructure:"prefix" yaml:"prefix"`
}

func NewRedisOptions() *RedisOptions {
	return &RedisOptions{
		Enabled: false,
		Address: "127.0.0.1:6379",
		Prefix:  "apiserver",
	}
}

// AddFlags adds flags related to mysql storage for a specific APIServer to the specified FlagSet.
func (o *RedisOptions) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&o.Enabled, "redis.enabled", o.Enabled, "Whether to enable Redis.")

	fs.StringVar(&o.Address, "redis.address", o.Address, ""+
		"Redis service host address. If left blank, the following related mysql options will be ignored.")

	fs.StringVar(&o.Prefix, "redis.prefix", o.Prefix, "Prefix identification key.")
}
