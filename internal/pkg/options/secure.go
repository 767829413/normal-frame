package options

import (
	"github.com/767829413/normal-frame/internal/pkg/config"
	"github.com/spf13/pflag"
)

// SecureOptions contains configuration items related to HTTPS or GRPC.
type SecureOptions struct {
	// ServerCert is the TLS cert info for serving secure traffic
	ServerCert GeneratableKeyCert `json:"tls" mapstructure:"tls"`
	// AdvertiseAddress net.IP
}

// NewSecureOptions creates a SecureOptions object with default parameters.
func NewSecureOptions() *SecureOptions {
	return &SecureOptions{
		ServerCert: GeneratableKeyCert{
			PairName:      "",
			CertDirectory: "",
		},
	}
}

// ApplyTo applies the run options to the method receiver and returns self.
func (s *SecureOptions) ApplyTo(ec *config.ExtraConfig) error {
	ec.CertDirectory = s.ServerCert.CertDirectory
	ec.PairName = s.ServerCert.PairName
	ec.CertFile = s.ServerCert.CertKey.CertFile
	ec.KeyFile = s.ServerCert.CertKey.KeyFile
	return nil
}

func (s *SecureOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.ServerCert.CertDirectory, "secure.tls.cert-dir", s.ServerCert.CertDirectory, ""+
		"The directory where the TLS certs are located. "+
		"If --secure.tls.cert-key.cert-file and --secure.tls.cert-key.private-key-file are provided, "+
		"this flag will be ignored.")

	fs.StringVar(&s.ServerCert.PairName, "secure.tls.pair-name", s.ServerCert.PairName, ""+
		"The name which will be used with --secure.tls.cert-dir to make a cert and key filenames. "+
		"It becomes <cert-dir>/<pair-name>.crt and <cert-dir>/<pair-name>.key")

	fs.StringVar(&s.ServerCert.CertKey.CertFile, "secure.tls.cert-key.cert-file", s.ServerCert.CertKey.CertFile, ""+
		"File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated "+
		"after server cert).")

	fs.StringVar(&s.ServerCert.CertKey.KeyFile, "secure.tls.cert-key.private-key-file",
		s.ServerCert.CertKey.KeyFile, ""+
			"File containing the default x509 private key matching --secure.tls.cert-key.cert-file.")
}

// GeneratableKeyCert contains configuration items related to certificate.
type GeneratableKeyCert struct {
	// CertKey allows setting an explicit cert/key file to use.
	CertKey CertKey `json:"cert-key" mapstructure:"cert-key"`

	// CertDirectory specifies a directory to write generated certificates to if CertFile/KeyFile aren't explicitly set.
	// PairName is used to determine the filenames within CertDirectory.
	// If CertDirectory and PairName are not set, an in-memory certificate will be generated.
	CertDirectory string `json:"cert-dir" mapstructure:"cert-dir"`
	// PairName is the name which will be used with CertDirectory to make a cert and key filenames.
	// It becomes CertDirectory/PairName.crt and CertDirectory/PairName.key
	PairName string `json:"pair-name" mapstructure:"pair-name"`
}

// CertKey contains configuration items related to certificate.
type CertKey struct {
	// CertFile is a file containing a PEM-encoded certificate, and possibly the complete certificate chain
	CertFile string `json:"cert-file" mapstructure:"cert-file"`
	// KeyFile is a file containing a PEM-encoded private key for the certificate specified by CertFile
	KeyFile string `json:"private-key-file" mapstructure:"private-key-file"`
}
