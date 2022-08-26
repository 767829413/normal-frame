package options

import (
	"github.com/767829413/normal-frame/internal/pkg/config"
	"github.com/gin-contrib/gzip"
	"github.com/spf13/pflag"
)

// FeatureOptions contains configuration items related to API server features.
type FeatureOptions struct {
	EnablePprof   bool `json:"enable-pprof" mapstructure:"enable-pprof" yaml:"enable-pprof"`
	EnableMetrics bool `json:"enable-metrics" mapstructure:"enable-metrics" yaml:"enable-metrics"`
	*Gzip
}

type Gzip struct {
	Enabled bool `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
	Level   int  `mapstructure:"level" json:"level" yaml:"level"`
}

// NewFeatureOptions creates a FeatureOptions object with default parameters.
func NewFeatureOptions() *FeatureOptions {
	return &FeatureOptions{
		EnableMetrics: false,
		EnablePprof:   false,
		Gzip: &Gzip{
			Enabled: false,
			Level:   gzip.DefaultCompression,
		},
	}
}

// ApplyTo applies the run options to the method receiver and returns self.
func (s *FeatureOptions) ApplyTo(c *config.GenericConfig) error {
	c.EnabledGzip = s.Gzip.Enabled
	c.GzipLevel = s.Gzip.Level
	c.EnableMetrics = s.EnableMetrics
	c.EnablePprof = s.EnablePprof
	return nil
}

// AddFlags adds flags for a specific APIServer to the specified FlagSet.
func (f *FeatureOptions) AddFlags(fs *pflag.FlagSet) {

	fs.BoolVar(&f.Gzip.Enabled, "feature.gzip.enabled", f.Gzip.Enabled, "Whether to enable Gzip compression.")

	fs.IntVar(&f.Gzip.Level, "feature.gzip.level", f.Gzip.Level, "The compression level can be any integer value between DefaultCompression = -1, NoCompression = 0, HuffmanOnly = -2 or BestSpeed = 1 and BestCompression = 9.")

	fs.BoolVar(&f.EnablePprof, "feature.enable-pprof", f.EnablePprof,
		"Enable pprof via web interface host:port/debug/pprof/")

	fs.BoolVar(&f.EnableMetrics, "feature.enable-metrics", f.EnableMetrics,
		"Enables metrics on the apiserver at /metrics")
}
