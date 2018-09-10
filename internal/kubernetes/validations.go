package kubernetes

import (
	"fmt"
	"strings"
)

type resolver interface {
	Resolve(string) (*Schema, error)
}

type Validator struct {
	resolver resolver
}

func NewValidator(resolver resolver) *Validator {
	return &Validator{
		resolver: resolver,
	}
}

// keep track of the key we're at key1.key2.key3....

func (v *Validator) Validate(incoming map[string]interface{}, schema *Schema, path []string) []error {
	errors := make([]error, 0)

	for key, value := range incoming {
		// Keep track of where we are
		tlp := make([]string, len(path))
		copy(tlp, path)
		tlp = append(tlp, key)

		property, ok := schema.Properties[key]
		if !ok {
			errors = append(errors, NewYamlPathError(tlp, NewUnknownKeyError(key)))
			continue
		}

		switch property.Type {
		case "string":
			if _, ok := value.(string); !ok {
				errors = append(errors, NewYamlPathError(tlp, NewWrongTypeError(key, "string", value)))
			}
		case "integer":
			// ignore property.Format until it causes a bug
			if _, ok := value.(int); !ok {
				errors = append(errors, NewYamlPathError(tlp, NewWrongTypeError(key, "integer", value)))
			}
		case "boolean":
			if _, ok := value.(bool); !ok {
				errors = append(errors, NewYamlPathError(tlp, NewWrongTypeError(key, "boolean", value)))
			}
		case "object":
			// this is for things like labels; map[interface{}]interface{} looks weird but that's how our yaml parser works.
			if _, ok := value.(map[interface{}]interface{}); !ok {
				errors = append(errors, NewYamlPathError(tlp, NewWrongTypeError(key, "map[interface{}]interface{}", value)))
			}
		case "array":
			items, ok := value.([]interface{})
			if !ok {
				errors = append(errors, NewYamlPathError(tlp, NewWrongTypeError(key, "[]interface{}", value)))
				continue
			}
			// TODO: check that items is not nil
			schema, err := v.resolver.Resolve(property.Items.Reference)
			if err != nil {
				fmt.Println(key, property)
				errors = append(errors, NewYamlPathError(tlp, err))
				continue
			}

			for _, item := range items {
				m, ok := item.(map[interface{}]interface{})
				if !ok {
					errors = append(errors, NewYamlPathError(tlp, NewWrongTypeError(key, "map[interface{}]interface{}", item)))
					continue
				}
				converted, err := keysToStrings(m)
				if err != nil {
					errors = append(errors, NewYamlPathError(tlp, err))
					continue
				}
				if errs := v.Validate(converted, schema, tlp); len(errs) > 0 {
					errors = append(errors, errs...)
					continue
				}
			}
		// default is some k8s object
		default:
			schema, err := v.resolver.Resolve(property.Reference)
			if err != nil {
				// DEBUG LINE good to use if there is a weird error
				// fmt.Println(key, property)
				errors = append(errors, NewYamlPathError(tlp, err))
				continue
			}
			d, ok := value.(map[interface{}]interface{})
			if !ok {
				errors = append(errors, NewYamlPathError(tlp, NewWrongTypeError(key, "map[interface{}]interface{}", value)))
				continue
			}
			convertedMap, err := keysToStrings(d)
			if err != nil {
				errors = append(errors, NewYamlPathError(tlp, err))
				continue
			}
			if subErrors := v.Validate(convertedMap, schema, tlp); len(subErrors) > 0 {
				errors = append(errors, subErrors...)
				continue
			}
		}
	}

	return errors
}

func keysToStrings(in map[interface{}]interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	for k, value := range in {
		key, ok := k.(string)
		if !ok {
			return nil, NewKeyNotStringError(k)
		}
		out[key] = value
	}
	return out, nil
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
	err  error
	path string
}

func NewYamlPathError(path []string, err error) error {
	return &YamlPathError{
		err:  err,
		path: strings.Join(path, "."),
	}
}
func (y *YamlPathError) Error() string {
	return fmt.Sprintf("[%s] %v", y.path, y.err)
}
