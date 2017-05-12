package envcnf

import "testing"

func Test_Parser_NewParser_Valid(t *testing.T) {
	var b bool
	p, err := NewParser(&b, "foo", "_")
	if err != nil {
		t.Fatalf("NewParser: %s", err)
	}
	t.Log("Parser:", p)
}

func Test_Parser_NewParser_ValueNotAPointer(t *testing.T) {
	var b bool
	p, err := NewParser(b, "foo", "_")
	if err == nil {
		t.Fatal("NewParser didn't error when receiving an non pointer value")
	}
	t.Log("Parser:", p)
}

func Test_Parser_NewParserWithName(t *testing.T) {
	var b bool
	p, err := NewParserWithName(&b, "foo", "_", "bar")
	if err != nil {
		t.Fatalf("NewParser: %s", err)
	}
	t.Log("Parser:", p)
}
