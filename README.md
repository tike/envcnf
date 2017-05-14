# envcnf

Parse values/settings from environment variables into arbitrarily nested/complex
golang data structures of any type with a single function call.

It is built to streamline the usage of [direnv](https://github.com/direnv/direnv)
for local/development needs and unix startup scripts for your webservice in
the production environment, while keeping your time/coding footprints small.
:rocket:

## The problem

 Passing configuration values to a daemon, e.g. a webserver, via environement
 variables can be tedious with go's onboard means and available libraries are
 limited in the types they support and in their handling of composite
 types, like structs, map and/or slices and especially arbitrarily nested
 combinations of these.

## The solution
 Package envcnf allows you to pass your configuration to your application via
 environment variables comfortably, reliably and without (practical) limitations.

## How

 1. Name your env vars to match your config type's field names. You can of
 course use a common prefix to set them apart from "regular" env vars and
 thous even have different configuration sets loaded at the same time,
 designating each set by a different prefix.

 2. Use [direnv](https:github.com/direnv/direnv) to maintain and organise
 your configuration in shell script like files or just source them by hand.

 3. Parse the configuration from the environment variables into your config
 struct, map, slice or any (however deeply nested and complex) combination
 thereof with **1 line of code**. See example right below ;-)

## Example:

Let the following be the data structs that hold your app's config:
```
package config

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
```

Your configuration file might then look like this:
```
# put this into your `.envrc`, if you're using `direnv`
# or add an invocation of your exe at the end and use it as a run script
# or `source` it into your shell and run your exe manually

# config values...
export ACME-COORP_Environment=production

export ACME-COORP_Listen_internal_Addr=127.0.0.1:80
export ACME-COORP_Listen_internal_HTTPS=false

export ACME-COORP_Listen_public_Addr=1.2.3.4:443
export ACME-COORP_Listen_public_HTTPS=true

export ACME-COORP_ChRoot=/var/empty

export ACME-COORP_MyFoo_Values_0=3
export ACME-COORP_MyFoo_Values_1=2
export ACME-COORP_MyFoo_Values_2=1
export ACME-COORP_MyFoo_Values_3=0

# to make this a run script, rather than a mere config, just add this one line
# /path/to/myapp
```

Now how do you parse those env vars into those structs?
```
package main

import (
  "github.com/tike/envcnf"

  "path/to/config"
)

func main(){
  
  var cnf config.MyCnf
  if err := envcnf.Parse(&cnf, "ACME-COORP", "_"); err != nil {
    fmt.Println("parsing config:", err)
    return
  }
  fmt.Printf("%#v\n",cnf)
}
```

Output will look like this:
```
config.MyCnf{
  Environment:"production",
  Listen:map[string]config.NetCnf{
    "internal":config.NetCnf{
      Addr:"127.0.0.1:80",
      HTTPS:false
    },
    "public":config.NetCnf{
      Addr:"1.2.3.4:443",
      HTTPS:true
    }
  },
  ChRoot:"/var/empty",
  MyFoo:config.MySection{
    Values:[]uint64{0x3, 0x2, 0x1, 0x0}
  }
}
```
