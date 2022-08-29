package options

import (
	"github.com/spf13/pflag"
)

type ApmOptions struct {
	Enabled bool   `json:"enabled" mapstructure:"enabled" yaml:"enabled"`
	Address string `mapstructure:"address" json:"address" yaml:"address"`
	Http    bool   `mapstructure:"http" json:"http" yaml:"http"`
	Mysql   bool   `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis   bool   `mapstructure:"redis" json:"redis" yaml:"redis"`
}

func NewApmOptions() *ApmOptions {
	return &ApmOptions{
		Enabled: true,
		Address: "/sidecar/sky-agent.sock",
		Http:    false,
		Mysql:   false,
		Redis:   false,
	}
}

func (o *ApmOptions) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&o.Enabled, "apm.enabled", o.Enabled, "Whether to enable APM.")

	fs.StringVar(&o.Address, "apm.address", o.Address, ""+
		"APM enabled address.")

	fs.BoolVar(&o.Http, "apm.http", o.Http, "Whether to enable Http.")

	fs.BoolVar(&o.Mysql, "apm.mysql", o.Mysql, "Whether to enable Mysql.")

	fs.BoolVar(&o.Redis, "apm.redis", o.Redis, "Whether to enable Redis.")

}
