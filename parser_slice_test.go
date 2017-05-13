package envcnf

import (
	"fmt"
	"os"
	"testing"
)

type TestSlice []int

func (ts TestSlice) setupEnv(t *testing.T, prefix, sepchar, name string) {
	basekey := ""
	if prefix != "" {
		basekey = prefix + sepchar
	}

	for k, v := range ts {
		os.Setenv(basekey+fmt.Sprintf("%s%s%v", name, sepchar, k), fmt.Sprintf("%v", v))
	}
}

func (ts TestSlice) teardownEnv(t *testing.T, prefix, sepchar, name string) {
	basekey := ""
	if prefix != "" {
		basekey = prefix + sepchar
	}

	for k := range ts {
		os.Unsetenv(basekey + fmt.Sprintf("%s%s%v", name, sepchar, k))
	}
}

var ts = TestSlice{11, -22, 33}

func Test_Parser_parseSlice_Valid_WithPrefix(t *testing.T) {
	ts.setupEnv(t, "ACME", "_", "SLICE")
	defer ts.teardownEnv(t, "ACME", "_", "SLICE")

	v := make(TestSlice, 0)
	p, err := NewParserWithName(&v, "ACME", "_", "SLICE")
	if err != nil {
		t.Fatalf("newParser: %#v", err)
	}

	if err := p.parseSlice(); err != nil {
		t.Fatalf("parseSlice said: %#v", err)
	}
	for k, v := range v {
		if v != ts[k] {
			t.Fatalf("failed to recover value\nHAVE: %#v\nWANT:%#v\n", v, ts)
		}
	}
}

func Test_Parser_parseSlice_Valid_WithoutPrefix(t *testing.T) {
	ts.setupEnv(t, "", "_", "SLICE")
	defer ts.teardownEnv(t, "", "_", "SLICE")

	v := make(TestSlice, 0)
	p, err := NewParserWithName(&v, "", "_", "SLICE")
	if err != nil {
		t.Fatalf("newParser: %#v", err)
	}
	if err := p.parseSlice(); err != nil {
		t.Fatalf("parseSlice said: %#v", err)
	}
	for k, v := range v {
		if v != ts[k] {
			t.Fatalf("failed to recover value\nHAVE: %#v\nWANT:%#v\n", v, ts)
		}
	}
}

func Test_Parser_parseSlice_InValid(t *testing.T) {
	var v TestSlice
	p, err := NewParserWithName(&v, "", "_", "ACME_FOO_THIS_VAR_SHOULD_NOT_EXIST")
	if err != nil {
		t.Fatalf("newParser: %#v", err)
	}
	if err := p.parseSlice(); err == nil {
		t.Fatal("parseSlice didn't error on non existing env var", err)
	}
}

type TestStringSlice []string

func (ts TestStringSlice) setupEnv(t *testing.T, prefix, sepchar, name string) {
	basekey := ""
	if prefix != "" {
		basekey = prefix + sepchar
	}

	for k, v := range ts {
		os.Setenv(basekey+fmt.Sprintf("%s%s%v", name, sepchar, k), fmt.Sprintf("%v", v))
	}
}

func (ts TestStringSlice) teardownEnv(t *testing.T, prefix, sepchar, name string) {
	basekey := ""
	if prefix != "" {
		basekey = prefix + sepchar
	}

	for k := range ts {
		os.Unsetenv(basekey + fmt.Sprintf("%s%s%v", name, sepchar, k))
	}
}

var tss = TestStringSlice{"a", "b", "c"}

func Test_Parser_parseSlice_string_Valid_WithPrefix(t *testing.T) {
	tss.setupEnv(t, "ACME", "_", "SLICE")
	defer tss.teardownEnv(t, "ACME", "_", "SLICE")

	v := make(TestStringSlice, 0)
	p, err := NewParserWithName(&v, "ACME", "_", "SLICE")
	if err != nil {
		t.Fatalf("newParser: %#v", err)
	}

	if err := p.parseSlice(); err != nil {
		t.Fatalf("parseSlice said: %#v", err)
	}
	for k, v := range v {
		if v != tss[k] {
			t.Fatalf("failed to recover value\nHAVE: %#v\nWANT:%#v\n", v, tss)
		}
	}
}

func Test_Parser_parseSlice_string_Valid_WithoutPrefix(t *testing.T) {
	tss.setupEnv(t, "", "_", "SLICE")
	defer tss.teardownEnv(t, "", "_", "SLICE")

	v := make(TestStringSlice, 0)
	p, err := NewParserWithName(&v, "", "_", "SLICE")
	if err != nil {
		t.Fatalf("newParser: %#v", err)
	}
	if err := p.parseSlice(); err != nil {
		t.Fatalf("parseSlice said: %#v", err)
	}
	for k, v := range v {
		if v != tss[k] {
			t.Fatalf("failed to recover value\nHAVE: %#v\nWANT:%#v\n", v, tss)
		}
	}
}
