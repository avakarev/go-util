package envutil

import "fmt"

const (
	// Var is env variable name
	Var = "APP_ENV"
	// Dev is "dev" env
	Dev = "dev"
	// Prod is "prod" env
	Prod = "prod"
)

// AppEnv implements application environment
type AppEnv string

// String returns env's string representation
func (env AppEnv) String() string {
	return string(env)
}

// IsDev checks whether given env is dev
func (env AppEnv) IsDev() bool {
	return env.String() == Dev
}

// IsProd checks whether given env is prod
func (env AppEnv) IsProd() bool {
	return env.String() == Prod
}

// NewAppEnv returns new AppEnv value
func NewAppEnv() (AppEnv, error) {
	s, err := MustStr(Var)
	if err != nil {
		return AppEnv(""), err
	}
	if s != Dev && s != Prod {
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
