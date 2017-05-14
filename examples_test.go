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
	os.Setenv("ACME-COORP_Environment", "production")

	os.Setenv("ACME-COORP_Listen_internal_Addr", "127.0.0.1:80")
	os.Setenv("ACME-COORP_Listen_internal_HTTPS", "false")

	os.Setenv("ACME-COORP_Listen_public_Addr", "1.2.3.4:443")
	os.Setenv("ACME-COORP_Listen_public_HTTPS", "true")

	os.Setenv("ACME-COORP_ChRoot", "/var/empty")

	os.Setenv("ACME-COORP_MyFoo_Values_0", "3")
	os.Setenv("ACME-COORP_MyFoo_Values_1", "2")
	os.Setenv("ACME-COORP_MyFoo_Values_2", "1")
	os.Setenv("ACME-COORP_MyFoo_Values_3", "0")

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
	if err := Parse(&config, "ACME-COORP", "_"); err != nil {
		// Handle error
	}

	fmt.Printf("%#v", config)
	/* gives the following output:
	   envcnf.MyCnf{
	     Environment:"production",
	     Listen:map[string]envcnf.NetCnf{
	       "internal":envcnf.NetCnf{
	         Addr:"127.0.0.1:80",
	         HTTPS:false
	       },
	       "public":envcnf.NetCnf{
	         Addr:"1.2.3.4:443",
	         HTTPS:true
	       }
	     },
	     ChRoot:"/var/empty",
	     MyFoo:envcnf.MySection{
	       Values:[]uint64{0x3, 0x2, 0x1, 0x0}
	     }
	   }
	*/
	// Output:
	// envcnf.MyCnf{Environment:"production", Listen:map[string]envcnf.NetCnf{"public":envcnf.NetCnf{Addr:"1.2.3.4:443", HTTPS:true}, "internal":envcnf.NetCnf{Addr:"127.0.0.1:80", HTTPS:false}}, ChRoot:"/var/empty", MyFoo:envcnf.MySection{Values:[]uint64{0x3, 0x2, 0x1, 0x0}}}
}

func TestExample(t *testing.T) {
	os.Setenv("ACME-COORP_Environment", "production")
	defer os.Unsetenv("ACME-COORP_Environment")

	os.Setenv("ACME-COORP_Listen_internal_Addr", "127.0.0.1:80")
	defer os.Unsetenv("ACME-COORP_Listen_internal_Addr")

	os.Setenv("ACME-COORP_Listen_internal_HTTPS", "false")
	defer os.Unsetenv("ACME-COORP_Listen_internal_HTTPS")

	os.Setenv("ACME-COORP_Listen_public_Addr", "1.2.3.4:443")
	defer os.Unsetenv("ACME-COORP_Listen_public_Addr")

	os.Setenv("ACME-COORP_Listen_public_HTTPS", "true")
	defer os.Unsetenv("ACME-COORP_Listen_public_HTTPS")

	os.Setenv("ACME-COORP_ChRoot", "/var/empty")
	defer os.Unsetenv("ACME-COORP_ChRoot")

	os.Setenv("ACME-COORP_MyFoo_Values_0", "3")
	defer os.Unsetenv("ACME-COORP_MyFoo_Values_0")

	os.Setenv("ACME-COORP_MyFoo_Values_1", "2")
	defer os.Unsetenv("ACME-COORP_MyFoo_Values_1")

	os.Setenv("ACME-COORP_MyFoo_Values_2", "1")
	defer os.Unsetenv("ACME-COORP_MyFoo_Values_2")

	os.Setenv("ACME-COORP_MyFoo_Values_3", "0")
	defer os.Unsetenv("ACME-COORP_MyFoo_Values_3")

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
		MyFoo: MySection{
			Values: []uint64{
				3, 2, 1, 0,
			},
		},
	}

	var config MyCnf
	if err := Parse(&config, "ACME-COORP", "_"); err != nil {
		t.Fatalf("Parse error: %v\n%#v\n", err, config)
	}
	if !reflect.DeepEqual(config, expect) {
		t.Fatalf("Unexpected Values parsed:\nHAVE:%#v\nWANT:%#v\n", config, expect)
	}
	t.Logf("parsed:\n%#v\n", config)
}
