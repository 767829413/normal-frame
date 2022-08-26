package config

type GenericConfig struct {
	Mode          string
	Healthz       bool
	BindAddress   string
	BindPort      int
	EnabledGzip   bool
	GzipLevel     int
	EnableMetrics bool
	EnablePprof   bool
}

// NewConfig returns a Config struct with the default values.
func NewGenericConfig() *GenericConfig {
	return &GenericConfig{}
}
