package envcnf

import (
	"errors"
	"fmt"
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

// parseString obtains the value from the env var that is signified by the fully
// nested (and possibly prefixed) name of the parser,
// parses it via strconv.ParseBool and assigns
// the obtained result to the (proper subfield) of the variable you handed to
// NewParser et al.
func (p *Parser) parseString() error {
	rawval, ok := p.env[p.getfullname()]
	if !ok {
		//TODO: use/obtain/signal default falue
		return errors.New("couldn't find envvar")
	}

	// CanAddr/CanSet/AssignableTo/ConvertibleTo are handled by the upper layers
	p.val.SetString(rawval)
	return nil
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

// parseInt obtains the value from the env var that is signified by the fully
// nested (and possibly prefixed) name of the parser,
// parses it via strconv.ParseBool and assigns
// the obtained result to the (proper subfield) of the variable you handed to
// NewParser et al.
func (p *Parser) parseInt() error {
	rawval, ok := p.env[p.getfullname()]
	if !ok {
		//TODO: use/obtain/signal default falue
		return errors.New("couldn't find envvar")
	}

	val, err := strconv.ParseInt(rawval, 10, p.valT.Bits())
	if err != nil {
		return err
	}

	// CanAddr/CanSet/AssignableTo/ConvertibleTo are handled by the upper layers
	p.val.SetInt(val)
	return nil
}

// parseUint obtains the value from the env var that is signified by the fully
// nested (and possibly prefixed) name of the parser,
// parses it via strconv.ParseBool and assigns
// the obtained result to the (proper subfield) of the variable you handed to
// NewParser et al.
func (p *Parser) parseUint() error {
	rawval, ok := p.env[p.getfullname()]
	if !ok {
		//TODO: use/obtain/signal default falue
		return errors.New("couldn't find envvar")
	}

	val, err := strconv.ParseUint(rawval, 10, p.valT.Bits())
	if err != nil {
		return err
	}

	// CanAddr/CanSet/AssignableTo/ConvertibleTo are handled by the upper layers
	p.val.SetUint(val)
	return nil
}

// parseFloat obtains the value from the env var that is signified by the fully
// nested (and possibly prefixed) name of the parser,
// parses it via strconv.ParseBool and assigns
// the obtained result to the (proper subfield) of the variable you handed to
// NewParser et al.
func (p *Parser) parseFloat() error {
	rawval, ok := p.env[p.getfullname()]
	if !ok {
		//TODO: use/obtain/signal default falue
		return errors.New("couldn't find envvar")
	}

	val, err := strconv.ParseFloat(rawval, p.valT.Bits())
	if err != nil {
		return err
	}

	// CanAddr/CanSet/AssignableTo/ConvertibleTo are handled by the upper layers
	p.val.SetFloat(val)
	return nil
}

// parseTypes invokes the correct handler method for the reflect.Kind of the
// value passed to NewParser et al.
func (p *Parser) parseTypes() error {
	switch p.val.Kind() {
	case reflect.Bool:
		return p.parseBool()
	case
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return p.parseInt()
	case
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		return p.parseUint()
	case
		reflect.Float32,
		reflect.Float64:
		return p.parseFloat()
	case reflect.Complex64, reflect.Complex128:
		return errors.New("complex parsing not yet implemented")
	case reflect.String:
		return p.parseString()
	case reflect.Ptr:
		return errors.New("ptr parsing not yet implemented")
	case reflect.Array, reflect.Slice:
		return errors.New("array/slice parsing not yet implemented")
	case reflect.Map:
		return errors.New("map parsing not yet implemented")
	case reflect.Struct:
		return p.parseStruct()
	default:
		return fmt.Errorf("parsing not yet implemented for value of Type %T (Kind: %s)", p.valT.Name(), p.valT.Kind())
	}
}

// parseStruct obtains the value from the env var that is signified by the fully
// nested (and possibly prefixed) name of the parser,
// parses it via strconv.ParseBool and assigns
// the obtained result to the (proper subfield) of the variable you handed to
// NewParser et al.
func (p *Parser) parseStruct() error {
	for i := 0; i < p.val.NumField(); i++ {
		field := p.val.Field(i)
		if !field.CanAddr() {
			return errors.New("struct field not addressable") //TODO: use dedicated error type with full info here
		}

		fieldName := p.valT.Field(i).Name
		subparser, err := newParserWithEnv(p.env, field.Addr().Interface(), p.prefix, p.sepchar, fieldName)
		if err != nil {
			return err
		}

		if len(p.parentNames) > 0 {
			subparser.parentNames = append(subparser.parentNames, p.parentNames...)
		}

		if field.Kind() == reflect.Struct {
			subparser.parentNames = append(subparser.parentNames, fieldName)
		}

		if err := subparser.parseTypes(); err != nil {
			return err
		}

	}
	return nil
}
