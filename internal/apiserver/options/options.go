// Package options contains flags and options for initializing an apiserver
package options

import (
	cliflag "github.com/767829413/normal-frame/fork/component-base/cli/flag"
	"github.com/767829413/normal-frame/internal/pkg/options"
)

type Options struct {
	GenericServerRunOptions *options.ServerRunOptions `json:"server" mapstructure:"server" yaml:"server"`
	MySQLOptions            *options.MySQLOptions     `json:"mysql" mapstructure:"mysql" yaml:"mysql"`
	RedisOptions            *options.RedisOptions     `json:"redis" mapstructure:"redis" yaml:"redis"`
	LogsOptions             *options.LogsOptions      `json:"logs" mapstructure:"logs" yaml:"logs"`
	GrpcOptions             *options.GrpcOptions      `json:"grpc" mapstructure:"grpc" yaml:"grpc"`
	FeatureOptions          *options.FeatureOptions   `json:"feature" mapstructure:"feature" yaml:"feature"`
	SecureOptions           *options.SecureOptions    `json:"secure" mapstructure:"secure" yaml:"secure"`
	HttpsOptions            *options.HttpsOptions     `json:"https" mapstructure:"https" yaml:"https"`
	ApmOptions              *options.ApmOptions       `json:"apm" mapstructure:"apm" yaml:"apm"`
}

// NewOptions creates a new Options object with default parameters.
func NewOptions() *Options {
	return &Options{
		GenericServerRunOptions: options.NewServerRunOptions(),
		MySQLOptions:            options.NewMySQLOptions(),
		RedisOptions:            options.NewRedisOptions(),
		LogsOptions:             options.NewLogsOptions(),
		GrpcOptions:             options.NewGrpcOptions(),
		FeatureOptions:          options.NewFeatureOptions(),
		SecureOptions:           options.NewSecureOptions(),
		HttpsOptions:            options.NewHttpsOptions(),
		ApmOptions:              options.NewApmOptions(),
	}
}

// Flags returns flags for a specific APIServer by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	o.GenericServerRunOptions.AddFlags(fss.FlagSet("server"))
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	o.LogsOptions.AddFlags(fss.FlagSet("logs"))
	o.GrpcOptions.AddFlags(fss.FlagSet("grpc"))
	o.FeatureOptions.AddFlags(fss.FlagSet("feature"))
	o.SecureOptions.AddFlags(fss.FlagSet("secure"))
	o.HttpsOptions.AddFlags(fss.FlagSet("https"))
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	o.ApmOptions.AddFlags(fss.FlagSet("apm"))
	return fss
}
