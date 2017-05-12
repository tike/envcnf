package envcnf

import (
	"errors"
	"reflect"
)

// Parser handles a single parsing process for a given (composite) value,
// thous allowing low overhead recursion to account for parsing of composite
// types.
type Parser struct {
	env rawEnv

	val  reflect.Value
	valT reflect.Type

	prefix  string
	sepchar string

	parentNames []string
	name        string
}

// NewParser is the default interface to be used for parsing composite types.
func NewParser(val interface{}, prefix, sepchar string) (*Parser, error) {
	env := newRawEnvWithPrfxSep(prefix, sepchar)
	return newParserWithEnv(env, val, prefix, sepchar, "")
}

// NewParserWithName is the default interface to be used for parsing a single
// non-composite value.
func NewParserWithName(val interface{}, prefix, sepchar, name string) (*Parser, error) {
	env := newRawEnvWithPrfxSep(prefix, sepchar)
	return newParserWithEnv(env, val, prefix, sepchar, name)
}

// newParserWithEnv constructs a Parser from the given values
func newParserWithEnv(env rawEnv, val interface{}, prefix, sepchar, name string) (*Parser, error) {
	ptrRef := reflect.ValueOf(val)
	if ptrRef.Kind() != reflect.Ptr {
		return nil, errors.New("val needs to be a pointer")
	}
	v := ptrRef.Elem()
	return &Parser{
		env: env,

		val:  v,
		valT: v.Type(),

		prefix:  prefix,
		sepchar: sepchar,
		name:    name,
	}, nil
}
