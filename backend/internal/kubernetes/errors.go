package kubernetes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// RequiredKeyNotFoundError is an error representing a required schema key that is missing
type RequiredKeyNotFoundError struct {
	key  string
	path []string
}

// NewRequiredKeyNotFoundError returns a new RequiredKeyNotFoundError
func NewRequiredKeyNotFoundError(key string, path []string) error {
	return &RequiredKeyNotFoundError{key: key, path: path}
}

// MarshalJSON extracts the error since error is an interface
func (r *RequiredKeyNotFoundError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Key   string
		Error string
	}{
		Key:   strings.Join(r.path, "."),
		Error: "Missing required key: " + r.key,
	})
}

// Error implements the error interface
func (r *RequiredKeyNotFoundError) Error() string {
	return fmt.Sprintf("key %q not found", r.key)
}

// UnknownTypeError represents an unknown type in a value.
// Examples are when the schema expects a string but gets an integer.
type UnknownTypeError struct {
	t interface{}
}

// NewUnknownTypeError returns a new UnknownTypeError
func NewUnknownTypeError(t interface{}) error {
	return &UnknownTypeError{t: t}
}

// Error implements the error interface
func (u *UnknownTypeError) Error() string {
	return fmt.Sprintf("unknown type %T", u.t)
}

// UnknownSchemaError means the resolver could not find the requested schema.
type UnknownSchemaError struct {
	schemaKey string
}

// NewUnknownSchemaError returns an UnknownSchemaError.
func NewUnknownSchemaError(key string) error {
	return &UnknownSchemaError{schemaKey: key}
}

// Error implements error interface.
func (u *UnknownSchemaError) Error() string {
	return fmt.Sprintf("unknown schema %v", u.schemaKey)
}

// MarshalJSON is used for the json.Marshal call.
func (u *UnknownSchemaError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Error string
		Key   string
	}{
		Error: u.Error(),
		Key:   "apiVersion",
	})
}

type KeyNotStringError struct {
	key interface{}
}

func NewKeyNotStringError(key interface{}) error {
	return &KeyNotStringError{key: key}
}
func (k *KeyNotStringError) Error() string {
	return fmt.Sprintf("key %v is a '%T' but needs to be a string", k.key, k.key)
}

type UnknownKeyError struct {
	key string
}

func NewUnknownKeyError(key string) error {
	return &UnknownKeyError{key: key}
}
func (u *UnknownKeyError) Error() string {
	return fmt.Sprintf("unknown key: %v", u.key)
}

type WrongTypeError struct {
	key          string
	expectedType string
	value        interface{}
}

func NewWrongTypeError(k, expectedType string, v interface{}) error {
	return &WrongTypeError{
		key:          k,
		expectedType: expectedType,
		value:        v,
	}
}
func (w *WrongTypeError) Error() string {
	return fmt.Sprintf("key %v has wrong type %T (should be %s)", w.key, w.value, w.expectedType)
}

type YamlPathError struct {
	Err   error
	Value interface{} `json:",omitempty"`
	Path  string
}

func NewYamlPathError(path []string, value interface{}, err error) error {
	return &YamlPathError{
		Err:   err,
		Value: value,
		Path:  strings.Join(path, "."),
	}
}
func (y *YamlPathError) Error() string {
	return fmt.Sprintf("[%s] %v", y.Path, y.Err)
}

// MarshalJSON extracts the error since error is an interface
func (y *YamlPathError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Key   string
		Value interface{}
		Error string
	}{
		Key:   y.Path,
		Value: y.Value,
		Error: y.Err.Error(),
	})
}

type UnknownFormatError struct {
	Format string
}

func (u *UnknownFormatError) Error() string {
	return fmt.Sprintf("unknown format %q", u.Format)
}
func NewUnknownFormatError(format string) error {
	return &UnknownFormatError{
		Format: format,
	}
}
