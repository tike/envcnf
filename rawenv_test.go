package envcnf

import (
	"os"
	"strings"
	"testing"
)

func Test_rawEnv_withprefix(t *testing.T) {
	should := rawEnv{
		"ACME_FOO": "1",
		"ACME_BAR": "2",
	}
	for k, v := range should {
		os.Setenv(k, v)
		defer os.Unsetenv(k)
	}

	is := newRawEnv("ACME_")
	for k, v := range should {
		isV, ok := is[strings.TrimPrefix(k, "ACME_")]
		if !ok {
			t.Errorf("Key %s not found", k)
		}
		if isV != v {
			t.Errorf("Value doesn't match: %s(should be: %s)", isV, v)
		}
	}
}

func Test_rawEnv_noprefix(t *testing.T) {
	env := newRawEnv("")
	if len(env) == 0 {
		t.Error("no environment variables found at all!")
	}
}

func Test_rawEnv_newRawEnvWithPrfxSep(t *testing.T) {
	env := newRawEnvWithPrfxSep("", "_")
	if len(env) == 0 {
		t.Error("no environment variables found at all!")
	}
}
