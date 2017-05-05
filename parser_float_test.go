package envcnf

import (
	"os"
	"testing"
)

func Test_Parser_parseFloat_Valid_WithPrefix(t *testing.T) {
	os.Setenv("ACME_FLOAT", "123.45")
	defer os.Unsetenv("ACME_FLOAT")

	var v float64
	p, err := NewParserWithName(&v, "ACME", "_", "FLOAT", NoConv)
	if err != nil {
		t.Fatalf("newParser: %#v", err)
	}
	if err := p.parseFloat(); err != nil {
		t.Fatalf("parseFloat said: %#v", err)
	}
	if v != 123.45 {
		t.Fatalf("failed to recover value")
	}
}

func Test_Parser_parseFloat_Valid_WithoutPrefix(t *testing.T) {
	os.Setenv("FLOAT", "123.45")
	defer os.Unsetenv("FLOAT")

	var v float64
	p, err := NewParserWithName(&v, "", "_", "FLOAT", NoConv)
	if err != nil {
		t.Fatalf("newParser: %#v", err)
	}
	if err := p.parseFloat(); err != nil {
		t.Fatalf("parseFloat said: %#v", err)
	}
	if v != 123.45 {
		t.Fatalf("failed to recover value")
	}
}

func Test_Parser_parseFloat_InValid(t *testing.T) {
	var v float64
	p, err := NewParserWithName(&v, "", "_", "ACME_FOO_THIS_VAR_SHOULD_NOT_EXIST", NoConv)
	if err != nil {
		t.Fatalf("newParser: %#v", err)
	}
	if err := p.parseFloat(); err == nil {
		t.Fatal("parseFloat didn't error on non existing env var", err)
	}
}
