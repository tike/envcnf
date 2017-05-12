package envcnf

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

type Inner struct {
	InnerA string
	InnerB uint64
}

type TestStruct struct {
	A string
	B int32
	C uint
	D float64
	E bool
	F Inner
}

func (ts TestStruct) setupEnv(t *testing.T, prefix, sepchar string) {
	basekey := ""
	if prefix != "" {
		basekey = prefix + sepchar
	}

	os.Setenv(basekey+"A", ts.A)
	os.Setenv(basekey+"B", fmt.Sprintf("%v", ts.B))
	os.Setenv(basekey+"C", fmt.Sprintf("%v", ts.C))
	os.Setenv(basekey+"D", fmt.Sprintf("%v", ts.D))
	os.Setenv(basekey+"E", fmt.Sprintf("%v", ts.E))
	os.Setenv(basekey+"F"+sepchar+"InnerA", fmt.Sprintf("%v", ts.F.InnerA))
	os.Setenv(basekey+"F"+sepchar+"InnerB", fmt.Sprintf("%v", ts.F.InnerB))
}

func (ts TestStruct) teardownEnv(t *testing.T, prefix, sepchar string) {
	basekey := ""
	if prefix != "" {
		basekey = prefix + sepchar
	}

	os.Unsetenv(basekey + "A")
	os.Unsetenv(basekey + "B")
	os.Unsetenv(basekey + "C")
	os.Unsetenv(basekey + "D")
	os.Unsetenv(basekey + "E")
	os.Unsetenv(basekey + "F" + sepchar + "InnerA")
	os.Unsetenv(basekey + "F" + sepchar + "InnerB")
}

var tc = TestStruct{
	A: "TestStructstring",
	B: -1,
	C: 500,
	D: -1.234456e+78,
	E: true,
	F: Inner{
		InnerA: "InnerA-Value",
		InnerB: 555666777,
	},
}

func Test_Parser_parseStruct_Valid_WithPrefix(t *testing.T) {
	tc.setupEnv(t, "ACME", "_")
	defer tc.teardownEnv(t, "ACME", "_")

	var v TestStruct
	p, err := NewParser(&v, "ACME", "_")
	if err != nil {
		t.Fatalf("newParser: %#v", err)
	}

	if err := p.parseStruct(); err != nil {
		t.Fatalf("parseStruct said: %#v", err)
	}
	if !reflect.DeepEqual(v, tc) {
		t.Fatalf("failed to recover value\nHAVE: %#v\nWANT:%#v\n", v, tc)
	}
}

func Test_Parser_parseStruct_Valid_WithoutPrefix(t *testing.T) {
	tc.setupEnv(t, "", "_")
	defer tc.teardownEnv(t, "", "_")

	var v TestStruct
	p, err := NewParserWithName(&v, "", "_", "STRUCT")
	if err != nil {
		t.Fatalf("newParser: %#v", err)
	}
	if err := p.parseStruct(); err != nil {
		t.Fatalf("parseStruct said: %#v", err)
	}
	if !reflect.DeepEqual(v, tc) {
		t.Fatalf("failed to recover value\nHAVE: %#v\nWANT:%#v\n", v, tc)
	}
}

func Test_Parser_parseStruct_InValid(t *testing.T) {
	var v TestStruct
	p, err := NewParserWithName(&v, "", "_", "ACME_FOO_THIS_VAR_SHOULD_NOT_EXIST")
	if err != nil {
		t.Fatalf("newParser: %#v", err)
	}
	if err := p.parseStruct(); err == nil {
		t.Fatal("parseStruct didn't error on non existing env var", err)
	}
}
