// Package envcnf parses envvars into composite types, e.g. for configuration.
//
// The problem:
//
// Passing configuration values to a daemon, e.g. a webserver, via environement
// variables can be tedious with go's onboard means and available libraries are
// limited in the types they support and in their handling of composite
// types, like structs, map and/or slices and especially arbitrarily nested
// combinations of these.
//
// The solution:
//
// Package envcnf allows you to pass your configuration to your application via
// environment variables comfortably, reliably and without (practical) limitations.
//
// How:
//
// 1. Name your env vars to match your config type's field names,
// case conversion to match all lower or all upper case is supported.
// You can of course use a common prefix to set them apart from
// "regular" env vars and thous even have different configuration sets loaded at
// the same time, designating each set by a different prefix.
//
// 2. Use direnv (https://github.com/direnv/direnv) to maintain and organise
// your configuration in shell script like files or just source them by hand.
//
// 3. Parse the configuration from the environment variables into your config
// struct, map, slice or any (however deeply nested and complex) combination
// thereof with __1 line of code__. Take a look at the examples ;-)
//
// 4. That's it! Spend your time on writing code that matters :-)
package envcnf
