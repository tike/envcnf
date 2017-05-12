package envcnf

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

type TestMap map[string]int

func (ts TestMap) setupEnv(t *testing.T, prefix, sepchar, name string) {
	basekey := ""
	if prefix != "" {
		basekey = prefix + sepchar
	}

	for k, v := range ts {
		os.Setenv(basekey+fmt.Sprintf("%s%s%v", name, sepchar, k), fmt.Sprintf("%v", v))
	}
}

func (ts TestMap) teardownEnv(t *testing.T, prefix, sepchar, name string) {
	basekey := ""
	if prefix != "" {
		basekey = prefix + sepchar
	}

	for k := range ts {
		os.Unsetenv(basekey + fmt.Sprintf("%s%s%v", name, sepchar, k))
	}
}

var tm = TestMap{
	"A": 11,
	"B": -22,
	"C": 33,
}

func Test_Parser_parseMap_Valid_WithPrefix(t *testing.T) {
	tm.setupEnv(t, "ACME", "_", "MAP")
	defer tm.teardownEnv(t, "ACME", "_", "MAP")

	v := make(TestMap)
	p, err := NewParserWithName(&v, "ACME", "_", "MAP")
	if err != nil {
		t.Fatalf("newParser: %#v", err)
	}

	if err := p.parseMap(); err != nil {
		t.Fatalf("parseMap said: %#v", err)
	}
	for k := range tm {
		v, ok := v[k]
		if !ok || v != tm[k] {
			t.Fatalf("failed to recover value\nHAVE: %#v\nWANT:%#v\n", v, tm)
		}
	}
}

func Test_Parser_parseMap_Valid_WithoutPrefix(t *testing.T) {
	tm.setupEnv(t, "", "_", "MAP")
	defer tm.teardownEnv(t, "", "_", "MAP")

	v := make(TestMap)
	p, err := NewParserWithName(&v, "", "_", "MAP")
	if err != nil {
		t.Fatalf("newParser: %#v", err)
	}
	if err := p.parseMap(); err != nil {
		t.Fatalf("parseMap said: %#v", err)
	}
	if !reflect.DeepEqual(v, tm) {
		t.Fatalf("failed to recover value\nHAVE: %#v\nWANT:%#v\n", v, tm)
	}
}

func Test_Parser_parseMap_InValid(t *testing.T) {
	var v TestMap
	p, err := NewParserWithName(&v, "", "_", "ACME_FOO_THIS_VAR_SHOULD_NOT_EXIST")
	if err != nil {
		t.Fatalf("newParser: %#v", err)
	}
	if err := p.parseMap(); err == nil {
		t.Fatal("parseMap didn't error on non existing env var", err)
	}
}
