package envutil

import "fmt"

const (
	// EnvDev is "dev" env
	EnvDev = "dev"
	// EnvBeta is "beta" env
	EnvBeta = "beta"
	// EnvProd is "prod" env
	EnvProd = "prod"
)

// AppEnv implements application environment
type AppEnv string

// String returns env's string representation
func (env AppEnv) String() string {
	return string(env)
}

// IsDev checks whether given env is dev
func (env AppEnv) IsDev() bool {
	return env.String() == EnvDev
}

// IsBeta checks whether given env is beta
func (env AppEnv) IsBeta() bool {
	return env.String() == EnvBeta
}

// IsProd checks whether given env is prod
func (env AppEnv) IsProd() bool {
	return env.String() == EnvProd
}

// NewAppEnv returns new AppEnv value
func NewAppEnv() (AppEnv, error) {
	s, err := MustStr("APP_ENV")
	if err != nil {
		return AppEnv(""), err
	}
	if s != EnvDev && s != EnvBeta && s != EnvProd {
		return AppEnv(""), fmt.Errorf("unexpected app env %q", s)
	}
	return AppEnv(s), nil
}

// MustAppEnv returns new AppEnv or panics if case of any errors
func MustAppEnv() AppEnv {
	env, err := NewAppEnv()
	if err != nil {
		panic(err)
	}
	return env
}
