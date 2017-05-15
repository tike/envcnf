package envcnf

import (
	"os"
	"testing"
)

type PointerInner struct {
	A int
}

type PointerTest struct {
	INT *PointerInner
}

func Test_Parser_parsePointer_Valid_WithPrefix(t *testing.T) {
	os.Setenv("ACME_INT_A", "123")
	defer os.Unsetenv("ACME_INT_A")

	var v PointerTest
	p, err := NewParser(&v, "ACME", "_", NoConv)
	if err != nil {
		t.Fatalf("newParser: %#v", err)
	}
	if err := p.Parse(); err != nil {
		t.Fatalf("parsePointer said: %T", err)
	}
	if v.INT.A != 123 {
		t.Fatalf("failed to recover value")
	}
}

func Test_Parser_parsePointer_Valid_WithoutPrefix(t *testing.T) {
	os.Setenv("INT_A", "123")
	defer os.Unsetenv("INT_A")

	var v PointerTest
	p, err := NewParser(&v, "", "_", NoConv)
	if err != nil {
		t.Fatalf("newParser: %#v", err)
	}
	if err := p.Parse(); err != nil {
		t.Fatalf("parsePointer said: %#v", err)
	}
	if v.INT.A != 123 {
		t.Fatalf("failed to recover value")
	}
}
