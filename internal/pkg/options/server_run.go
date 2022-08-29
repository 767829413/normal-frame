package options

import (
	"github.com/767829413/normal-frame/internal/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
)

type ServerRunOptions struct {
	Mode        string `json:"mode" mapstructure:"mode" yaml:"mode"`
	Healthz     bool   `json:"healthz" mapstructure:"healthz" yaml:"healthz"`
	BindAddress string `json:"bind-address" mapstructure:"bind-address" yaml:"bind-address"`
	BindPort    int    `json:"bind-port" mapstructure:"bind-port" yaml:"bind-port"`
}

func NewServerRunOptions() *ServerRunOptions {
	return &ServerRunOptions{
		Mode:        gin.ReleaseMode,
		Healthz:     true,
		BindAddress: "",
		BindPort:    80,
	}
}

// ApplyTo applies the run options to the method receiver and returns self.
func (s *ServerRunOptions) ApplyTo(c *config.GenericConfig) error {
	c.Mode = s.Mode
	c.Healthz = s.Healthz
	c.BindAddress = s.BindAddress
	c.BindPort = s.BindPort
	return nil
}

func (s *ServerRunOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.Mode, "server.mode", s.Mode, ""+
		"Start the server in a specified server mode. Supported server mode: debug, test, release.")

	fs.BoolVar(&s.Healthz, "server.healthz", s.Healthz, ""+
		"Add self readiness check and install /healthz router.")

	fs.StringVar(&s.BindAddress, "server.bind-address", s.BindAddress, ""+
		"The IP address on which to serve the --insecure.bind-port "+
		"(set to 0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces).")

	fs.IntVar(&s.BindPort, "server.bind-port", s.BindPort, "The port the service is listening on.")
}
