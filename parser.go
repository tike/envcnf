package envcnf

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
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

// getfullname concatenates the parts of the parsers (parent) name(s) in a
// sensible way.
func (p Parser) getfullname() string {
	var key string
	if len(p.parentNames) > 0 {
		key = strings.Join(p.parentNames, p.sepchar) + p.sepchar
	}
	return key + p.name
}

// parseBool obtains the value from the env var that is signified by the fully
// nested (and possibly prefixed) name of the parser,
// parses it via strconv.ParseBool and assigns
// the obtained result to the (proper subfield) of the variable you handed to
// NewParser et al.
func (p *Parser) parseBool() error {
	rawval, ok := p.env[p.getfullname()]
	if !ok {
		//TODO: use/obtain/signal default falue
		return errors.New("couldn't find envvar")
	}
	val, err := strconv.ParseBool(rawval)
	if err != nil {
		return err
	}
	// CanAddr/CanSet/AssignableTo/ConvertibleTo are handled by the upper layers
	p.val.SetBool(val)
	return nil
}
