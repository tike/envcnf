package envcnf

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

// Here's how you'd use the Parse function with a sample config.
func ExampleParse() {
	/* Here is how a sample config file could look like:

		    # put this into your .envrc, if you're using direnv
		    # or name it `run.sh` and add an invocation of your exe at the end
		    # or `source` it manually in your shell and run exe from there

		    # config values...
	      export ACME-COORP_Environment=production

				export ACME-COORP_Listen_internal_Addr=127.0.0.1:80
				export ACME-COORP_Listen_internal_HTTPS=false # or any of: FALSE, False, 1, f, F

				export ACME-COORP_Listen_public_Addr=1.2.3.4:443
				export ACME-COORP_Listen_public_HTTPS=true # or any of: TRUE, True, 1, t, T

			  export ACME-COORP_ChRoot=/var/empty

			  export ACME-COORP_MyFoo_Values_0=3
				export ACME-COORP_MyFoo_Values_1=2
				export ACME-COORP_MyFoo_Values_2=1
				export ACME-COORP_MyFoo_Values_3=0

		    # ./acme_secretsauce
	*/

	// meanwhile in your source code...
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

	var config MyCnf
	if err := Parse(&config, "ACME-COORP", "_"); err != nil {
		// Handle error
	}
	fmt.Println(config)
	/*
			    MyCnf{
			    Environment: "production",
			    Listen: map{
			    "internal": NetCnf{
			      Addr: "127.0.0.1:80",
			      HTTPS: false,
			    }
			    "external": NetCnf{
			      Addr: "1.2.3.4:443",
			      HTTPS: false,
			    }
		      ChRoot: "/var/empty",
		      MyFoo: int{3, 2, 1, 0},
			  }
	*/
	// TADA !
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
	config.Listen = make(map[string]NetCnf)
	config.MyFoo.Values = make([]uint64, 0)
	t.Logf("before parse: %#v\n", config)

	if err := Parse(&config, "ACME-COORP", "_"); err != nil {
		t.Fatalf("Parse: %v\n%v\n", config, err)
	}
	if !reflect.DeepEqual(config, expect) {
		t.Fatalf("Unexpected Values parsed:\nHAVE:%#v\nWANT:%#v\n", config, expect)
	}
}
