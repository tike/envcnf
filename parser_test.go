package envcnf

import (
	"os"
	"testing"
)

func Test_Parser_NewParser_Valid(t *testing.T) {
	var v bool
	p, err := NewParser(&v, "foo", "_", NoConv)
	if err != nil {
		t.Fatalf("NewParser: %s", err)
	}
	t.Log("Parser:", p)
}

func Test_Parser_NewParser_ValueNotAPointer(t *testing.T) {
	var v bool
	p, err := NewParser(v, "foo", "_", NoConv)
	if err == nil {
		t.Fatal("NewParser didn't error when receiving an non pointer value")
	}
	t.Log("Parser:", p)
}

func Test_Parser_NewParserWithName(t *testing.T) {
	var v bool
	p, err := NewParserWithName(&v, "foo", "_", "bar", NoConv)
	if err != nil {
		t.Fatalf("NewParser: %s", err)
	}
	t.Log("Parser:", p)
}

func Test_Parser_Parse(t *testing.T) {
	os.Setenv("ACME_INT", "123")
	defer os.Unsetenv("ACME_INT")

	var v int
	p, err := NewParserWithName(&v, "ACME", "_", "INT", NoConv)
	if err != nil {
		t.Fatalf("NewParser: %s", err)
	}
	if err := p.Parse(); err != nil {
		t.Fatalf("Parser.Parse(): %v", err)
	}
}

func Test_Parser_getfullname(t *testing.T) {
	var v bool
	p, err := NewParserWithName(&v, "foo", "_", "testname", NoConv)
	if err != nil {
		t.Fatalf("NewParser: %s", err)
	}
	name := p.getfullname()
	if name != "testname" {
		t.Fatalf("name is %s (expected: %s)", name, "testname")
	}
}

func Test_Parser_getfullname_complete(t *testing.T) {
	var v bool
	p, err := NewParserWithName(&v, "foo", "_", "testname", NoConv)
	if err != nil {
		t.Fatalf("NewParser: %s", err)
	}
	p.parentNames = []string{"a", "b", "c"}
	name := p.getfullname()
	if name != "a_b_c_testname" {
		t.Fatalf("name is %s (expected: %s)", name, "a_b_c_testname")
	}

}
