package options

import (
	"github.com/767829413/normal-frame/internal/pkg/config"
	"github.com/spf13/pflag"
)

// GRPCOptions are for creating an unauthenticated, unauthorized, insecure port.
// No one should be using these anymore.
type GrpcOptions struct {
	Enabled     bool   `json:"enabled" mapstructure:"enabled" yaml:"enabled"`
	BindAddress string `json:"bind-address" mapstructure:"bind-address"`
	BindPort    int    `json:"bind-port"    mapstructure:"bind-port"`
	MaxMsgSize  int    `json:"max-msg-size" mapstructure:"max-msg-size"`
}

// NewGRPCOptions is for creating an unauthenticated, unauthorized, insecure port.
// No one should be using these anymore.
func NewGrpcOptions() *GrpcOptions {
	return &GrpcOptions{
		BindAddress: "0.0.0.0",
		BindPort:    8081,
		MaxMsgSize:  4 * 1024 * 1024,
	}
}

func (s *GrpcOptions) ApplyTo(ec *config.ExtraConfig) error {
	ec.EnableGRPC = s.Enabled
	ec.GrpcAddress = s.BindAddress
	ec.GrpcPort = s.BindPort
	ec.MaxMsgSize = s.MaxMsgSize
	return nil
}

// AddFlags adds flags related to features for a specific api server to the
// specified FlagSet.
func (s *GrpcOptions) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&s.Enabled, "grpc.enabled", s.Enabled, "Whether to enable GRPC.")

	fs.StringVar(&s.BindAddress, "grpc.bind-address", s.BindAddress, ""+
		"The IP address on which to serve the --grpc.bind-port(set to 0.0.0.0 for all IPv4 interfaces and :: for all IPv6 interfaces).")

	fs.IntVar(&s.BindPort, "grpc.bind-port", s.BindPort, ""+
		"The port on which to serve unsecured, unauthenticated grpc access. It is assumed "+
		"that firewall rules are set up such that this port is not reachable from outside of "+
		"the deployed machine and that port 443 on the iam public address is proxied to this "+
		"port. This is performed by nginx in the default setup. Set to zero to disable.")

	fs.IntVar(&s.MaxMsgSize, "grpc.max-msg-size", s.MaxMsgSize, "gRPC max message size.")
}
