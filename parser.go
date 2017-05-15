package envcnf

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// These values are used to indicate wether to do case conversion when looking up
// environment variable names. So if a struct field is named 'Field'
// and you pass 'ToUpper' the parser will look for an environment variable
// named 'FIELD'.
const (
	NoConv int = iota
	ToLower
	ToUpper
)

// Parser handles a single parsing process for a given (composite) value,
// thous allowing low overhead recursion to account for parsing of composite
// types.
type Parser struct {
	env rawEnv

	val  reflect.Value
	valT reflect.Type

	conv    int
	prefix  string
	sepchar string

	parentNames []string
	name        string
}

// Parse is the main interface to the package. Just pass a pointer to the variable
// you'd like to receive your config values in. If you use a common prefix to
// set your config variable names apart and avoid cluttering, pass it via the
// prefix parameter. SepChar is used to separate the prefix and the subfields
// of your env var. See the examples.
func Parse(val interface{}, prefix, sepchar string) error {
	p, err := NewParser(val, prefix, sepchar)
	if err != nil {
		return err
	}
	return p.parseTypes()
}

// NewParser can be used parse multiple values into a composite types, like
// structs, maps or slices. Just pass a pointer to the variable
// you'd like to receive your config values in. If you use a common prefix to
// set your config variable names apart and avoid cluttering, pass it via the
// prefix parameter. sepchar is used to separate the prefix and the subfields
// of your env var. See the examples.
func NewParser(val interface{}, prefix, sepchar string) (*Parser, error) {
	env := newRawEnvWithPrfxSep(prefix, sepchar)
	return newParserWithEnv(env, val, prefix, sepchar, "")
}

// NewParserWithName can be used to parse a single non-composite value from an
// environment variable. Just pass a pointer to the variable
// you'd like to receive your config value in. If you use a common prefix to
// set your config variable names apart and avoid cluttering, pass it via the
// prefix parameter. sepchar is used to separate the prefix and the subfields
// of your env var. See the examples.
func NewParserWithName(val interface{}, prefix, sepchar, name string) (*Parser, error) {
	env := newRawEnvWithPrfxSep(prefix, sepchar)
	return newParserWithEnv(env, val, prefix, sepchar, name)
}

// newParserWithEnv constructs a Parser from the given values
func newParserWithEnv(env rawEnv, val interface{}, prefix, sepchar, name string) (*Parser, error) {
	ref := reflect.ValueOf(val)
	if ref.Kind() != reflect.Ptr && ref.Kind() != reflect.Interface {
		return nil, ErrNeedPointerValue
	}
	v := ref.Elem()
	return &Parser{
		env: env,

		val:  v,
		valT: v.Type(),

		prefix:  prefix,
		sepchar: sepchar,
		name:    name,
	}, nil
}

// Parse starts the parsing process, returning any errors encountered.
func (p *Parser) Parse() error {
	return p.parseTypes()
}

// getfullname concatenates the parts of the parser's (parent) name(s) in a
// sensible way.
func (p Parser) getfullname() string {
	var key string
	if len(p.parentNames) > 0 {
		key = strings.Join(p.parentNames, p.sepchar) + p.sepchar
	}
	key += p.name

	return p.convertCase(key)
}

func (p Parser) convertCase(key string) string {
	switch p.conv {
	case ToUpper:
		return strings.ToUpper(key)
	case ToLower:
		return strings.ToLower(key)
	default:
		return key
	}
}

// parseString obtains the value from the env var that is signified by the fully
// nested (and possibly prefixed) name of the parser,
// parses it via strconv.ParseString, expands any contained evironment variables
// and assigns the obtained result to the (proper subfield of the) variable you
// handed to NewParser or NewParserWithName.
func (p *Parser) parseString() error {
	key := p.getfullname()
	rawval, ok := p.env[key]
	if !ok {
		//TODO: use/obtain/signal default value
		return MissingEnvVar(key)
	}

	// this should almost never be necessary,
	// but it's nice to have.
	rawval = os.ExpandEnv(rawval)

	// CanAddr/CanSet/AssignableTo/ConvertibleTo are handled by the upper layers
	p.val.SetString(rawval)
	return nil
}

// parseBool obtains the value from the env var that is signified by the fully
// nested (and possibly prefixed) name of the parser,
// parses it via strconv.ParseBool and assigns
// the obtained result to the (proper subfield of the) variable you handed to
// NewParser or NewParserWithName.
func (p *Parser) parseBool() error {
	key := p.getfullname()
	rawval, ok := p.env[key]
	if !ok {
		//TODO: use/obtain/signal default value
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
// the obtained result to the (proper subfield of the) variable you handed to
// NewParser or NewParserWithName.
func (p *Parser) parseInt() error {
	key := p.getfullname()
	rawval, ok := p.env[key]
	if !ok {
		//TODO: use/obtain/signal default value
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
// the obtained result to the (proper subfield of the) variable you handed to
// NewParser or NewParserWithName.
func (p *Parser) parseUint() error {
	key := p.getfullname()
	rawval, ok := p.env[key]
	if !ok {
		//TODO: use/obtain/signal default value
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
// the obtained result to the (proper subfield of the) variable you handed to
// NewParser or NewParserWithName.
func (p *Parser) parseFloat() error {
	key := p.getfullname()
	rawval, ok := p.env[key]
	if !ok {
		//TODO: use/obtain/signal default value
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

func (p *Parser) parsePointer() error {
	if p.val.IsNil() {
		if !p.val.CanSet() {
			return FieldNotAddressable(p.name + " not settable")
		}
		p.val.Set(reflect.New(p.valT.Elem()))
	}

	subEnv := p.env.getAllWithPrefix(p.name + p.sepchar)
	subparser, err := newParserWithEnv(subEnv, p.val.Interface(), p.prefix, p.sepchar, "")
	if err != nil {
		return err
	}
	if len(p.parentNames) > 0 {
		subparser.parentNames = append(subparser.parentNames, p.parentNames...)
	}
	if p.val.Kind() == reflect.Struct {
		subparser.parentNames = append(subparser.parentNames, p.name)
	}

	return subparser.parseTypes()
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
		return p.parsePointer()
	case reflect.Array, reflect.Slice:
		return p.parseSlice()
	case reflect.Map:
		return p.parseMap()
	case reflect.Struct:
		return p.parseStruct()
	default:
		fmt.Println("unsupported:", p.valT, p.val.Interface())
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

		// TODO: This should not be necessary, so there's something odd here
		// or right above in the invocation of newParserWithEnv
		if field.Kind() == reflect.Slice {
			field.Set(subparser.val)
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
		return MissingEnvVar(prfx + p.sepchar + "KEY for map value")
	}

	if p.val.IsNil() {
		p.val.Set(reflect.MakeMap(p.valT))
	}

	keyT := p.valT.Key()
	keyIsString := keyT.Kind() == reflect.String

	valT := p.valT.Elem()
	var valIsString, valIsContainer bool
	switch valT.Kind() {
	case reflect.String:
		valIsString = true
	case reflect.Struct, reflect.Slice, reflect.Array, reflect.Map:
		valIsContainer = true
	}

	for k, v := range env {
		var mapKey, subTypeKey string

		convertedKey := reflect.New(keyT)
		if !keyIsString {
			valParser, err := newParserWithEnv(env, convertedKey.Interface(), "", p.sepchar, "")
			if err != nil {
				return err
			}
			if err := valParser.parseTypes(); err != nil {
				return err
			}
		} else if valIsContainer {
			parts := strings.SplitN(k, p.sepchar, 2)
			mapKey, subTypeKey = parts[0], parts[1]
		} else {
			mapKey, subTypeKey = k, k
		}
		convertedKey.Elem().SetString(mapKey)

		convertedVal := reflect.New(valT)
		if valIsString {
			convertedVal.Elem().SetString(v)
		} else if !valIsContainer {
			valParser, err := newParserWithEnv(env, convertedVal.Interface(), "", p.sepchar, subTypeKey)
			if err != nil {
				return err
			}
			if err := valParser.parseTypes(); err != nil {
				return err
			}
		} else {
			subEnv := env.getAllWithPrefix(mapKey + p.sepchar)
			for subK := range subEnv {
				valParser, err := newParserWithEnv(subEnv, convertedVal.Interface(), "", p.sepchar, subK)
				if err != nil {
					return err
				}
				if err := valParser.parseTypes(); err != nil {
					return err
				}
			}
		}
		p.val.SetMapIndex(convertedKey.Elem(), convertedVal.Elem())
	}
	return nil
}

// parseSlice obtains all values from the env vars that are prefixed by the fully
// nested (and possibly prefixed) name of the parser,
// parses them recursively and assigns
// the obtained result to the (proper subfield of the) variable you handed to
// NewParser or NewParserWithName.
func (p *Parser) parseSlice() error {
	prfx := p.getfullname()
	env := p.env.getAllWithPrefix(prfx + p.sepchar)

	if len(env) == 0 {
		return MissingEnvVar(prfx + p.sepchar + "N for slice/array value")
	}

	if p.val.IsNil() {
		p.val.Set(reflect.MakeSlice(p.valT, 0, len(env)))
	}

	valT := p.valT.Elem()
	var valIsString, valIsContainer bool
	switch valT.Kind() {
	case reflect.String:
		valIsString = true
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Struct:
		valIsContainer = true
	}

	convertedVals := make(map[int]reflect.Value)
	for k, v := range env {
		// parse index into int
		if valIsContainer {
			parts := strings.SplitN(k, p.sepchar, 2)
			k = parts[0]
		}

		idx, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			return err
		}

		// setup value type pointer
		convertedVal := reflect.New(valT)
		if valIsString {
			convertedVal.Elem().SetString(v)
		} else if !valIsContainer {
			valParser, err := newParserWithEnv(env, convertedVal.Interface(), "", "", k)
			if err != nil {
				return err
			}
			if err := valParser.parseTypes(); err != nil {
				return err
			}
		} else {
			if _, ok := convertedVals[int(idx)]; ok {
				continue
			}
			subEnv := env.getAllWithPrefix(k + p.sepchar)
			for subK := range subEnv {
				valParser, err := newParserWithEnv(subEnv, convertedVal.Interface(), "", p.sepchar, subK)
				if err != nil {
					return err
				}
				if err := valParser.parseTypes(); err != nil {
					return err
				}
			}
		}
		// collect unorderd
		convertedVals[int(idx)] = convertedVal
	}

	// finally add values to target container in designated order
	for i := 0; i < len(convertedVals); i++ {
		p.val = reflect.Append(p.val, convertedVals[i].Elem())
	}
	return nil
}
