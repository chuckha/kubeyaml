package kubernetes

import (
	"fmt"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Input struct {
	Kind       string
	APIVersion string
	Data       map[string]interface{}
}

type Loader struct{}

func NewLoader() *Loader {
	return &Loader{}
}

func (l *Loader) Load(reader io.Reader) (*Input, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read incoming reader: %v", err)
	}

	incoming := map[string]interface{}{}
	if err := yaml.Unmarshal(b, incoming); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml with error %v", err)
	}

	val, ok := incoming["apiVersion"]
	if !ok {
		return nil, NewRequiredKeyNotFoundError("apiVersion")
	}

	apiVersion, ok := val.(string)
	if !ok {
		return nil, NewUnknownTypeError(val)
	}

	val, ok = incoming["kind"]
	if !ok {
		return nil, NewRequiredKeyNotFoundError("kind")
	}

	kind, ok := val.(string)
	if !ok {
		return nil, NewUnknownTypeError(val)
	}

	delete(incoming, "apiVersion")
	delete(incoming, "kind")
	return &Input{
		APIVersion: apiVersion,
		Kind:       kind,
		Data:       incoming,
	}, nil
}

type RequiredKeyNotFoundError struct {
	key string
}

func NewRequiredKeyNotFoundError(key string) error {
	return &RequiredKeyNotFoundError{key: key}
}
func (r *RequiredKeyNotFoundError) Error() string {
	return fmt.Sprintf("key %q not found", r.key)
}

type UnknownTypeError struct {
	t interface{}
}

func NewUnknownTypeError(t interface{}) error {
	return &UnknownTypeError{t: t}
}
func (u *UnknownTypeError) Error() string {
	return fmt.Sprintf("unknown type %T", u.t)
}
