package options

import (
	"github.com/767829413/normal-frame/internal/pkg/config"
	"github.com/spf13/pflag"
)

// GRPCOptions are for creating an unauthenticated, unauthorized, insecure port.
// No one should be using these anymore.
type HttpsOptions struct {
	Enabled     bool   `json:"enabled" mapstructure:"enabled" yaml:"enabled"`
	BindAddress string `json:"bind-address" mapstructure:"bind-address"`
	BindPort    int    `json:"bind-port"    mapstructure:"bind-port"`
}

// NewHttpsOptions contains configuration items related to HTTPS server startup.
func NewHttpsOptions() *HttpsOptions {
	return &HttpsOptions{
		Enabled:     false,
		BindAddress: "0.0.0.0",
		BindPort:    8443,
	}
}

func (s *HttpsOptions) ApplyTo(ec *config.ExtraConfig) error {
	ec.EnableHttps = s.Enabled
	ec.HttpsAddress = s.BindAddress
	ec.HttpsPort = s.BindPort
	return nil
}

func (s *HttpsOptions) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&s.Enabled, "https.enabled", s.Enabled, "Whether to enable GRPC.")

	fs.StringVar(&s.BindAddress, "https.bind-address", s.BindAddress, ""+
		"The IP address on which to listen for the --https.bind-port port. The "+
		"associated interface(s) must be reachable by the rest of the engine, and by CLI/web "+
		"clients. If blank, all interfaces will be used (0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces).")

	fs.IntVar(&s.BindPort, "https.bind-port", s.BindPort, "The port on which to serve HTTPS with authentication and authorization.")
}
