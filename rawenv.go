package envcnf

import (
	"os"
	"strings"
)

// rawEnv is used to obtain and handle sets of environment variables.
type rawEnv map[string]string

// newRawEnv obtains all environment variables with the given prefix via os.Environ,
// and packages those in a rawEnv map type.
// If prefix is the empty string,
// all variables defined in the processes's environment are being returned.
func newRawEnv(prefix string) rawEnv {
	rawvals := os.Environ()

	var env rawEnv
	if prefix == "" {
		env = make(rawEnv, len(rawvals))
	} else {
		env = make(rawEnv)
	}

	for _, rawval := range rawvals {
		if !strings.HasPrefix(rawval, prefix) {
			continue
		}
		keyval := strings.SplitN(rawval, "=", 2)
		key := strings.TrimPrefix(keyval[0], prefix)
		env[key] = keyval[1]

	}
	return env
}

// newRawEnvWithPrfxSep returns the full env if prefix is the empty string,
// otherwise the limited subset of env vars that begin with prefix+sepchar
// is selected and prefix+sepchar is stripped from the env var names.
func newRawEnvWithPrfxSep(prefix, sepchar string) rawEnv {
	var env rawEnv
	if prefix == "" {
		env = newRawEnv(prefix)
	} else {
		env = newRawEnv(prefix + sepchar)
	}
	return env
}

// getAllWithPrefix returns the subset of values in the map that start with the
// given prefix, the prefix is stripped from the keys in the returned map.
func (r rawEnv) getAllWithPrefix(prefix string) rawEnv {
	sub := make(rawEnv)
	for k, v := range r {
		if strings.HasPrefix(k, prefix) {
			sub[strings.TrimPrefix(k, prefix)] = v
		}
	}
	return sub
}
