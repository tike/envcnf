package envcnf

import (
	"os"
	"testing"
)

func TestConvToUpper(t *testing.T) {
	os.Setenv("ACME_BOOL", "true")
	defer os.Unsetenv("ACME_BOOL")

	var v bool
	p, err := NewParserWithName(&v, "acme", "_", "bool", ToUpper)
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

func TestConvToLower(t *testing.T) {
	os.Setenv("acme_bool", "true")
	defer os.Unsetenv("acme_bool")

	var v bool
	p, err := NewParserWithName(&v, "ACME", "_", "BOOL", ToLower)
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
