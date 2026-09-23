package envutil_test

import (
	"errors"
	"testing"

	"github.com/avakarev/go-util/envutil"
	"github.com/avakarev/go-util/testutil"
)

func TestAppEnvString(t *testing.T) {
	testutil.Diff("dev", envutil.AppEnv("dev").String(), t)
}

func TestAppEnvIs(t *testing.T) {
	dev := envutil.AppEnv(envutil.EnvDev)
	testutil.Diff(true, dev.IsDev(), t)
	testutil.Diff(false, dev.IsBeta(), t)
	testutil.Diff(false, dev.IsProd(), t)

	beta := envutil.AppEnv(envutil.EnvBeta)
	testutil.Diff(false, beta.IsDev(), t)
	testutil.Diff(true, beta.IsBeta(), t)
	testutil.Diff(false, beta.IsProd(), t)

	prod := envutil.AppEnv(envutil.EnvProd)
	testutil.Diff(false, prod.IsDev(), t)
	testutil.Diff(false, prod.IsBeta(), t)
	testutil.Diff(true, prod.IsProd(), t)
}

func TestNewAppEnv(t *testing.T) {
	_, err := envutil.NewAppEnv()
	testutil.MustErr(errors.New("envutil: env variable \"APP_ENV\" is required, but not set"), err, t)

	resetEnv := testutil.SetEnv(testutil.Env{"APP_ENV": "staging"})
	_, err = envutil.NewAppEnv()
	testutil.MustErr(errors.New("unexpected app env \"staging\""), err, t)
	resetEnv()

	resetEnv = testutil.SetEnv(testutil.Env{"APP_ENV": "prod"})
	defer resetEnv()
	env, err := envutil.NewAppEnv()
	testutil.MustNoErr(err, t)
	testutil.Diff(envutil.AppEnv("prod"), env, t)
}

func TestMustAppEnv(t *testing.T) {
	func() {
		defer func() {
			testutil.Diff(true, recover() != nil, t)
		}()
		_ = envutil.MustAppEnv()
	}()

	resetEnv := testutil.SetEnv(testutil.Env{"APP_ENV": "beta"})
	defer resetEnv()
	testutil.Diff(envutil.AppEnv("beta"), envutil.MustAppEnv(), t)
}
