package mitmproxy

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"online/common/log"
	"online/common/utils"
	"online/common/utils/tlsutils"
	"online/mitmproxy/mitm"
	"time"
)

type Config struct {
	Host            string        `json:"host"`
	Port            int           `json:"port"`
	Ca              []byte        `json:"ca"`
	Key             []byte        `json:"key"`
	TransparentMode bool          `json:"transparent_mode"`
	Timeout         time.Duration `json:"timeout"`

	mitmConfig *mitm.Config

	webhookCallback        func(req *http.Request)
	mirrorRequestCallback  func(req *http.Request)
	mirrorResponseCallback func(req *http.Response)
}

type Option func(config *Config)

func WithMirrorRequest(cb func(req *http.Request)) Option {
	return func(config *Config) {
		config.mirrorRequestCallback = func(req *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Error(err)
				}
			}()
			cb(req)
		}
	}
}

func WithMirrorResponse(cb func(req *http.Response)) Option {
	return func(config *Config) {
		config.mirrorResponseCallback = func(req *http.Response) {
			defer func() {
				if err := recover(); err != nil {
					log.Error(err)
				}
			}()
			cb(req)
		}
	}
}

func WithWebHook(cb func(req *http.Request)) Option {
	return func(config *Config) {
		config.webhookCallback = func(req *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Error(err)
				}
			}()
			cb(req)
		}
	}
}

func WithHost(host string) Option {
	return func(config *Config) {
		config.Host = host
	}
}

func WithDefaultTimeout(d float64) Option {
	return func(config *Config) {
		config.Timeout = utils.FloatSecondDuration(d)
	}
}

func WithPort(port int) Option {
	return func(config *Config) {
		config.Port = port
	}
}

func WithCaCert(ca []byte, key []byte) Option {
	return func(config *Config) {
		config.Ca = ca
		config.Key = key
	}
}

func WithAutoCa() Option {
	return func(config *Config) {
		var err error
		config.Ca, config.Key, err = tlsutils.GenerateSelfSignedCertKeyWithCommonName("CA-for-MITM", "", nil, nil)
		if err != nil {
			log.Errorf("generate self signed cert failed: %s", err)
		}
	}
}

func WithTransparentMode(b bool) Option {
	return func(config *Config) {
		config.TransparentMode = b
	}
}

func NewConfig(opts ...Option) (*Config, error) {
	config := &Config{
		Host: "0.0.0.0", Port: 8088,
	}
	for _, opt := range opts {
		opt(config)
	}

	ca, key := config.Ca, config.Key
	if ca == nil || key == nil {
		return nil, utils.Error("empty ca-cert or key...")
	}

	c, err := tls.X509KeyPair(ca, key)
	if err != nil {
		return nil, utils.Errorf("parse ca and privKey failed: %s", err)
	}

	cert, err := x509.ParseCertificate(c.Certificate[0])
	if err != nil {
		return nil, utils.Errorf("extract x509 cert failed: %s", err)
	}

	mc, err := mitm.NewConfig(cert, c.PrivateKey)
	if err != nil {
		return nil, utils.Errorf("build private key failed: %s", err)
	}
	mc.SkipTLSVerify(true)
	mc.SetOrganization("MITMServer")
	mc.SetValidity(time.Hour * 24 * 365)
	config.mitmConfig = mc
	return config, nil
}
