package options

import "github.com/spf13/pflag"

type LogsOptions struct {
	OutPut string `json:"out-put" mapstructure:"out-put" yaml:"out-put"`
}

func NewLogsOptions() *LogsOptions {
	return &LogsOptions{
		OutPut: "redis",
	}
}

func (o *LogsOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.OutPut, "log.output", o.OutPut, "Log output path, stdout or redis, default redis.")
}
