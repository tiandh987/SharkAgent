package server

import (
	"github.com/gin-gonic/gin"
)

// Config 是用于配置 GenericAPIServer 的结构体。
type Config struct {
	Mode            string
	Middlewares     []string
	InsecureServing *InsecureServingInfo
	SecureServing   *SecureServingInfo

	//JWT *JwtInfo

	Healthz         bool
	EnableProfiling bool
	EnableMetrics   bool
}

// NewConfig returns a Config struct with the default values.
func NewConfig() *Config {
	return &Config{
		Mode:            gin.ReleaseMode,
		Middlewares:     []string{},
		Healthz:         true,
		EnableProfiling: true,
		EnableMetrics:   true,
		//Jwt: &JwtInfo{
		//	Realm:      "iam jwt",
		//	Timeout:    1 * time.Hour,
		//	MaxRefresh: 1 * time.Hour,
		//},
	}
}

// Complete fills in any fields not set that are required to have valid data and can be derived
// from other fields. If you're going to `ApplyOptions`, do that first. It's mutating the receiver.
func (c *Config) Complete() CompletedConfig {
	return CompletedConfig{c}
}

// ===============================================================

// InsecureServingInfo holds configuration of the insecure http server.
type InsecureServingInfo struct {
	Address string
}

// ================================================================

// CertKey contains configuration items related to certificate.
type CertKey struct {
	// CertFile 是一个包含 PEM 编码证书的文件，可能还有完整的证书链
	CertFile string
	// KeyFile 是一个文件，其中包含 CertFile 指定的证书的 PEM 编码私钥
	KeyFile string
}

// SecureServingInfo holds configuration of the TLS server.
type SecureServingInfo struct {
	BindAddress string
	BindPort    int
	CertKey     CertKey
}

// ================================================================

// CompletedConfig is the completed configuration for GenericAPIServer.
type CompletedConfig struct {
	*Config
}

// New returns a new instance of GenericAPIServer from the given config.
func (c CompletedConfig) New() (*GenericAPIServer, error) {
	s := &GenericAPIServer{
		SecureServingInfo:   c.SecureServing,
		InsecureServingInfo: c.InsecureServing,
		mode:                c.Mode,
		healthz:             c.Healthz,
		enableMetrics:       c.EnableMetrics,
		enableProfiling:     c.EnableProfiling,
		middlewares:         c.Middlewares,
		Engine:              gin.New(),
	}

	//initGenericAPIServer(s)

	return s, nil
}

