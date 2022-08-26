package config

// ExtraConfig defines extra configuration for the apiserver.
type ExtraConfig struct {
	EnableHttps   bool
	EnableGRPC    bool
	HttpsAddress  string
	HttpsPort     int
	GrpcAddress   string
	GrpcPort      int
	MaxMsgSize    int
	CertDirectory string
	PairName      string
	CertFile      string
	KeyFile       string
}

func NewExtraConfig() *ExtraConfig {
	return &ExtraConfig{}
}
