package envcnf

import (
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
		return nil, ErrNeedPointerValue
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
// NewParser or NewParserWithName.
func (p *Parser) parseString() error {
	key := p.getfullname()
	rawval, ok := p.env[key]
	if !ok {
		//TODO: use/obtain/signal default falue
		return MissingEnvVar(key)
	}

	// CanAddr/CanSet/AssignableTo/ConvertibleTo are handled by the upper layers
	p.val.SetString(rawval)
	return nil
}

// parseBool obtains the value from the env var that is signified by the fully
// nested (and possibly prefixed) name of the parser,
// parses it via strconv.ParseBool and assigns
// the obtained result to the (proper subfield) of the variable you handed to
// NewParser or NewParserWithName.
func (p *Parser) parseBool() error {
	key := p.getfullname()
	rawval, ok := p.env[key]
	if !ok {
		//TODO: use/obtain/signal default falue
		return MissingEnvVar(key)
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
// parses it via strconv.ParseInt and assigns
// the obtained result to the (proper subfield) of the variable you handed to
// NewParser or NewParserWithName.
func (p *Parser) parseInt() error {
	key := p.getfullname()
	rawval, ok := p.env[key]
	if !ok {
		//TODO: use/obtain/signal default falue
		return MissingEnvVar(key)
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
// parses it via strconv.ParseUint and assigns
// the obtained result to the (proper subfield) of the variable you handed to
// NewParser or NewParserWithName.
func (p *Parser) parseUint() error {
	key := p.getfullname()
	rawval, ok := p.env[key]
	if !ok {
		//TODO: use/obtain/signal default falue
		return MissingEnvVar(key)
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
// parses it via strconv.ParseFloat and assigns
// the obtained result to the (proper subfield) of the variable you handed to
// NewParser or NewParserWithName.
func (p *Parser) parseFloat() error {
	key := p.getfullname()
	rawval, ok := p.env[key]
	if !ok {
		//TODO: use/obtain/signal default falue
		return MissingEnvVar(key)
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
// value passed to NewParser or NewParserWithName.
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
		return UnsupportedType("Complex64/Complex128")
	case reflect.String:
		return p.parseString()
	case reflect.Ptr:
		return UnsupportedType("Ptr")
	case reflect.Array, reflect.Slice:
		return UnsupportedType("Array/Slice")
	case reflect.Map:
		return p.parseMap()
	case reflect.Struct:
		return p.parseStruct()
	default:
		return UnsupportedType(p.valT.Name() + " of kind " + p.valT.Kind().String())
	}
}

// parseStruct obtains the values from the env vars that are signified by the fully
// nested (and possibly prefixed) name of the parser,
// parses them recursively and assigns
// the obtained result to the (proper subfield of the) variable you handed to
// NewParser or NewParserWithName.
func (p *Parser) parseStruct() error {
	for i := 0; i < p.val.NumField(); i++ {
		field := p.val.Field(i)
		fieldName := p.valT.Field(i).Name
		if !field.CanAddr() {
			return FieldNotAddressable(fieldName)
		}

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

// parseMap obtains all values from the env vars that are prefixed by the fully
// nested (and possibly prefixed) name of the parser,
// parses them recursively and assigns
// the obtained result to the (proper subfield of the) variable you handed to
// NewParser or NewParserWithName.
func (p *Parser) parseMap() error {
	prfx := p.getfullname()
	env := p.env.getAllWithPrefix(prfx + p.sepchar)

	if len(env) == 0 {
		return MissingEnvVar(prfx + "_XYZ for map value")
	}

	keyT := p.valT.Key()
	needKeyTrans := keyT.Kind() != reflect.String

	valT := p.valT.Elem()
	needValTrans := valT.Kind() != reflect.String

	for k, v := range env {
		convertedKey := reflect.New(keyT)
		if needKeyTrans {
			valParser, err := newParserWithEnv(env, convertedKey.Interface(), "", "", "")
			if err != nil {
				return err
			}
			if err := valParser.parseTypes(); err != nil {
				return err
			}
		} else {
			convertedKey.Elem().SetString(k)
		}

		convertedVal := reflect.New(valT)
		if needValTrans {
			valParser, err := newParserWithEnv(env, convertedVal.Interface(), "", "", k)
			if err != nil {
				return err
			}
			if err := valParser.parseTypes(); err != nil {
				return err
			}
		} else {
			convertedVal.Elem().SetString(v)
		}
		p.val.SetMapIndex(convertedKey.Elem(), convertedVal.Elem())

	}
	return nil
}
