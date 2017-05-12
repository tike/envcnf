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

func Test_rawEnv_getAllWithPrefix(t *testing.T) {
	env := rawEnv{
		"aa": "c",
		"ab": "c",
		"ac": "c",
		"b":  "b",
		"ca": "b",
	}
	res := env.getAllWithPrefix("a")
	if len(res) != 3 {
		t.Errorf("not all valid values selected: have %d (expected: %d)", len(res), 3)
	}
	for k, v := range res {
		if v != "c" {
			t.Errorf("wrong key selected: %s", k)
		}
	}
}
