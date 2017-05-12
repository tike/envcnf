package envcnf

import (
	"errors"
	"fmt"
)

// ErrNeedPointerValue is returned by NewParser et al if it is called with a plain
// value instead of a pointer.
var ErrNeedPointerValue = errors.New("envcnf: val needs to be a pointer")

// MissingEnvVar is returned when no env var with a name fitting the scheme
// for given field can be found.
type MissingEnvVar string

func (e MissingEnvVar) Error() string {
	return fmt.Sprintf("envcnf: missing env var %q", e)
}

// FieldNotAddressable is returned when a structfield is not addressable.
// see the reflect.Addr and reflect.CanAddr if you don't know what this means.
type FieldNotAddressable string

func (e FieldNotAddressable) Error() string {
	return fmt.Sprintf("envcnf: struct field %q not addressable", e)
}

// UnsupportedType is returned when a type is not supported.
type UnsupportedType string

func (e UnsupportedType) Error() string {
	return fmt.Sprintf("envcnf: unsupported type %q", e)
}
