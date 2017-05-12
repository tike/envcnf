package envcnf

import (
	"os"
	"testing"
)

func Test_Parser_parseUint_Valid_WithPrefix(t *testing.T) {
	os.Setenv("ACME_UINT", "123")
	defer os.Unsetenv("ACME_UINT")

	var v uint
	p, err := NewParserWithName(&v, "ACME", "_", "UINT")
	if err != nil {
		t.Fatalf("newParser: %#v", err)
	}
	if err := p.parseUint(); err != nil {
		t.Fatalf("parseUint said: %#v", err)
	}
	if v != 123 {
		t.Fatalf("failed to recover value")
	}
}

func Test_Parser_parseUint_Valid_WithoutPrefix(t *testing.T) {
	os.Setenv("UINT", "123")
	defer os.Unsetenv("UINT")

	var v uint
	p, err := NewParserWithName(&v, "", "_", "UINT")
	if err != nil {
		t.Fatalf("newParser: %#v", err)
	}
	if err := p.parseUint(); err != nil {
		t.Fatalf("parseUint said: %#v", err)
	}
	if v != 123 {
		t.Fatalf("failed to recover value")
	}
}

func Test_Parser_parseUint_InValid(t *testing.T) {
	var v uint
	p, err := NewParserWithName(&v, "", "_", "ACME_FOO_THIS_VAR_SHOULD_NOT_EXIST")
	if err != nil {
		t.Fatalf("newParser: %#v", err)
	}
	if err := p.parseUint(); err == nil {
		t.Fatal("parseUint didn't error on non existing env var", err)
	}
}
