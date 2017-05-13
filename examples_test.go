package envcnf

import (
	"fmt"
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
	//TODO: copy stuff from above...
}
