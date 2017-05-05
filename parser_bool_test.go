package envcnf

import (
	"os"
	"testing"
)

func Test_Parser_parseBool_Valid_WithPrefix(t *testing.T) {
	os.Setenv("ACME_BOOL", "true")
	defer os.Unsetenv("ACME_BOOL")

	var v bool
	p, err := NewParser(&v, "ACME", "_", NoConv)
	if err != nil {
		t.Fatalf("newParser: %#v", err)
	}
	p.name = "BOOL"
	if err := p.parseBool(); err != nil {
		t.Fatalf("parseBool said: %#v", err)
	}
	if !v {
		t.Fatalf("failed to recover value")
	}
}

func Test_Parser_parseBool_Valid_WithoutPrefix(t *testing.T) {
	os.Setenv("BOOL", "true")
	defer os.Unsetenv("BOOL")

	var v bool
	p, err := NewParserWithName(&v, "", "_", "BOOL", NoConv)
	if err != nil {
		t.Fatalf("newParser: %#v", err)
	}
	if err := p.parseBool(); err != nil {
		t.Fatalf("parseBool said: %#v", err)
	}
	if !v {
		t.Fatalf("failed to recover value")
	}
}

func Test_Parser_parseBool_InValid(t *testing.T) {
	var v bool
	p, err := NewParserWithName(&v, "", "_", "ACME_FOO_THIS_VAR_SHOULD_NOT_EXIST", NoConv)
	if err != nil {
		t.Fatalf("newParser: %#v", err)
	}
	if err := p.parseBool(); err == nil {
		t.Fatal("parseBool didn't error on non existing env var", err)
	}
}
