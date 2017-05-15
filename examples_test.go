package envcnf

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

// Here's how you'd use the Parse function with a sample config.
func ExampleParse() {
	// you'd normally set those outside the program, of course.
	os.Setenv("ACME-CORP_Environment", "production")

	os.Setenv("ACME-CORP_Listen_internal_Addr", "127.0.0.1:80")
	os.Setenv("ACME-CORP_Listen_internal_HTTPS", "false")

	os.Setenv("ACME-CORP_Listen_public_Addr", "1.2.3.4:443")
	os.Setenv("ACME-CORP_Listen_public_HTTPS", "true")

	os.Setenv("ACME-CORP_ChRoot", "/var/empty")

	os.Setenv("ACME-CORP_MyFoo_Values_0", "3")
	os.Setenv("ACME-CORP_MyFoo_Values_1", "2")
	os.Setenv("ACME-CORP_MyFoo_Values_2", "1")
	os.Setenv("ACME-CORP_MyFoo_Values_3", "0")

	// here are the sample configuration types we'll be parsing into
	type NetCnf struct {
		Addr  string
		HTTPS bool
	}

	type MySection struct {
		Values []uint64
	}

	type MyCnf struct {
		Environment string
		Listen      map[string]NetCnf
		ChRoot      string
		MyFoo       MySection
	}

	// and this is how you'd parse this
	var config MyCnf
	if err := Parse(&config, "ACME-CORP", "_"); err != nil {
		fmt.Println("Parse:", err)
		return
	}

	fmt.Println(config)
	// Output should look like:
	// {production map[internal:{127.0.0.1:80 false} public:{1.2.3.4:443 true}] /var/empty {[3 2 1 0]}}
}

func ExampleNewParser() {
	// you'd normally set those outside the program, of course.
	os.Setenv("ACME-CORP_Host", "localhost")
	os.Setenv("ACME-CORP_Port", "8000")
	os.Setenv("ACME-CORP_HTTPS", "false")

	type Addr struct {
		Host  string
		Port  int
		HTTPS bool
	}

	var cnf Addr
	p, err := NewParser(&cnf, "ACME-CORP", "_")
	if err != nil {
		fmt.Println("NewParser:", err)
		return
	}

	if err := p.Parse(); err != nil {
		fmt.Println("Parse:", err)
		return
	}

	fmt.Println(cnf)
	// Output:
	// {localhost 8000 false}
}

func ExampleNewParserWithName() {
	// you'd normally set those outside the program, of course.
	os.Setenv("ACME-CORP_environment", "production")

	var env string
	p, err := NewParserWithName(&env, "ACME-CORP", "_", "environment")
	if err != nil {
		fmt.Println("NewParserWithName:", err)
		return
	}

	if err := p.Parse(); err != nil {
		fmt.Println("Parse:", err)
		return
	}

	fmt.Println(env)

	// Output:
	// production
}

func TestFooBar(t *testing.T) {
	os.Setenv("ACME-CORP_Environment", "production")
	defer os.Unsetenv("ACME-CORP_Environment")

	os.Setenv("ACME-CORP_Listen_internal_Addr", "127.0.0.1:80")
	defer os.Unsetenv("ACME-CORP_Listen_internal_Addr")

	os.Setenv("ACME-CORP_Listen_internal_HTTPS", "false")
	defer os.Unsetenv("ACME-CORP_Listen_internal_HTTPS")

	os.Setenv("ACME-CORP_Listen_public_Addr", "1.2.3.4:443")
	defer os.Unsetenv("ACME-CORP_Listen_public_Addr")

	os.Setenv("ACME-CORP_Listen_public_HTTPS", "true")
	defer os.Unsetenv("ACME-CORP_Listen_public_HTTPS")

	os.Setenv("ACME-CORP_ChRoot", "/var/empty")
	defer os.Unsetenv("ACME-CORP_ChRoot")

	os.Setenv("ACME-CORP_MyFoo_Values_0", "3")
	defer os.Unsetenv("ACME-CORP_MyFoo_Values_0")

	os.Setenv("ACME-CORP_MyFoo_Values_1", "2")
	defer os.Unsetenv("ACME-CORP_MyFoo_Values_1")

	os.Setenv("ACME-CORP_MyFoo_Values_2", "1")
	defer os.Unsetenv("ACME-CORP_MyFoo_Values_2")

	os.Setenv("ACME-CORP_MyFoo_Values_3", "0")
	defer os.Unsetenv("ACME-CORP_MyFoo_Values_3")

	type NetCnf struct {
		Addr  string
		HTTPS bool
	}

	type MySection struct {
		Values []uint64
	}

	type MyCnf struct {
		Environment string
		Listen      map[string]NetCnf
		ChRoot      string
		MyFoo       *MySection
	}

	var expect = MyCnf{
		Environment: "production",
		Listen: map[string]NetCnf{
			"internal": NetCnf{
				Addr:  "127.0.0.1:80",
				HTTPS: false,
			},
			"public": NetCnf{
				Addr:  "1.2.3.4:443",
				HTTPS: true,
			},
		},
		ChRoot: "/var/empty",
		MyFoo: &MySection{
			Values: []uint64{
				3, 2, 1, 0,
			},
		},
	}

	var config MyCnf
	if err := Parse(&config, "ACME-CORP", "_"); err != nil {
		t.Fatalf("Parse error: %v\n%#v\n", err, config)
	}

	if config.Environment != expect.Environment {
		t.Fatalf("Unexpected Values parsed:\nHAVE:%#v\nWANT:%#v\n", config.Environment, expect.Environment)
	}
	if config.ChRoot != expect.ChRoot {
		t.Fatalf("Unexpected Values parsed:\nHAVE:%#v\nWANT:%#v\n", config.ChRoot, expect.ChRoot)
	}
	if !reflect.DeepEqual(config.Listen["internal"], expect.Listen["internal"]) {
		t.Fatalf("Unexpected Values parsed:\nHAVE:%#v\nWANT:%#v\n", config.Listen["internal"], expect.Listen["internal"])
	}
	if !reflect.DeepEqual(config.Listen["public"], expect.Listen["public"]) {
		t.Fatalf("Unexpected Values parsed:\nHAVE:%#v\nWANT:%#v\n", config.Listen["public"], expect.Listen["public"])
	}
	if !reflect.DeepEqual(*config.MyFoo, *expect.MyFoo) {
		t.Fatalf("Unexpected Values parsed:\nHAVE:%#v\nWANT:%#v\n", *config.MyFoo, *expect.MyFoo)
	}

	t.Logf("parsed:\n%#v\n", config)
}
