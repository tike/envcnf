package envcnf

import (
	"os"
	"testing"
)

func Test_Parser_parseString_Valid_WithPrefix(t *testing.T) {
	os.Setenv("ACME_STRING", "foomatic2000")
	defer os.Unsetenv("ACME_STRING")

	var v string
	p, err := NewParser(&v, "ACME", "_")
	if err != nil {
		t.Fatalf("newParser: %#v", err)
	}
	p.name = "STRING"
	if err := p.parseString(); err != nil {
		t.Fatalf("parseString said: %#v", err)
	}
	if v != "foomatic2000" {
		t.Fatalf("failed to recover value")
	}
}

func Test_Parser_parseString_Valid_WithoutPrefix(t *testing.T) {
	os.Setenv("STRING", "foomatic2000")
	defer os.Unsetenv("STRING")

	var v string
	p, err := NewParserWithName(&v, "", "_", "STRING")
	if err != nil {
		t.Fatalf("newParser: %#v", err)
	}
	if err := p.parseString(); err != nil {
		t.Fatalf("parseString said: %#v", err)
	}
	if v != "foomatic2000" {
		t.Fatalf("failed to recover value")
	}
}

func Test_Parser_parseString_InValid(t *testing.T) {
	var v string
	p, err := NewParserWithName(&v, "", "_", "ACME_FOO_THIS_VAR_SHOULD_NOT_EXIST")
	if err != nil {
		t.Fatalf("newParser: %#v", err)
	}
	if err := p.parseString(); err == nil {
		t.Fatal("parseString didn't error on non existing env var", err)
	}
}
