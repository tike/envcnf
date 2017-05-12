package envcnf

import (
	"os"
	"testing"
)

func Test_Parser_parseInt_Valid_WithPrefix(t *testing.T) {
	os.Setenv("ACME_INT", "123")
	defer os.Unsetenv("ACME_INT")

	var v int
	p, err := NewParserWithName(&v, "ACME", "_", "INT")
	if err != nil {
		t.Fatalf("newParser: %#v", err)
	}
	if err := p.parseInt(); err != nil {
		t.Fatalf("parseInt said: %#v", err)
	}
	if v != 123 {
		t.Fatalf("failed to recover value")
	}
}

func Test_Parser_parseInt_Valid_WithoutPrefix(t *testing.T) {
	os.Setenv("INT", "123")
	defer os.Unsetenv("INT")

	var v int
	p, err := NewParserWithName(&v, "", "_", "INT")
	if err != nil {
		t.Fatalf("newParser: %#v", err)
	}
	if err := p.parseInt(); err != nil {
		t.Fatalf("parseInt said: %#v", err)
	}
	if v != 123 {
		t.Fatalf("failed to recover value")
	}
}

func Test_Parser_parseInt_InValid(t *testing.T) {
	var v int
	p, err := NewParserWithName(&v, "", "_", "ACME_FOO_THIS_VAR_SHOULD_NOT_EXIST")
	if err != nil {
		t.Fatalf("newParser: %#v", err)
	}
	if err := p.parseInt(); err == nil {
		t.Fatal("parseInt didn't error on non existing env var", err)
	}
}
