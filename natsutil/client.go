package natsutil

import (
	"time"

	"github.com/avakarev/go-util/envutil"
)

// ClientConfig defines client config
type ClientConfig struct {
	Env     string
	Timeout time.Duration
}

// Copy returns config copy
func (config *ClientConfig) Copy() *ClientConfig {
	return &ClientConfig{
		Env:     config.Env,
		Timeout: config.Timeout,
	}
}

// ClientOption defines client configuring func type
type ClientOption func(config *ClientConfig)

// ClientConfigure applies given options to given config
func ClientConfigure(config *ClientConfig, opts []ClientOption) *ClientConfig {
	if len(opts) == 0 {
		return config
	}
	cp := config.Copy()
	for _, fn := range opts {
		fn(cp)
	}
	return cp
}

// WithEnv sets env to given value
func WithEnv(env string) ClientOption {
	return func(config *ClientConfig) {
		config.Env = env
	}
}

// WithTimeout sets timeout to given value
func WithTimeout(timeout time.Duration) ClientOption {
	return func(config *ClientConfig) {
		config.Timeout = timeout
	}
}

// DefaultClientConfig returns config with default values
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		Env:     envutil.EnvProd,
		Timeout: 8 * time.Second,
	}
}
