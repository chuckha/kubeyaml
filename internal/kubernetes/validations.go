package kubernetes

import "fmt"

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

func (v *Validator) Validate(incoming map[string]interface{}, schema *Schema) []error {
	errors := make([]error, 0)
	// key validation
	for key, value := range incoming {
		property, ok := schema.Properties[key]
		if !ok {
			errors = append(errors, NewUnknownKeyError(key))
			continue
		}

		switch property.Type {
		case "string":
			if _, ok := value.(string); !ok {
				errors = append(errors, NewWrongTypeError(key, "string", value))
				continue
			}
		// not a built in type
		default:
			schema, err := v.resolver.Resolve(property.Reference)
			if err != nil {
				errors = append(errors, err)
				continue
			}
			d, ok := value.(map[string]interface{})
			if !ok {
				errors = append(errors, NewWrongTypeError(key, "map[string]interface{}", value))
				continue
			}
			if subErrors := v.Validate(d, schema); len(subErrors) > 0 {
				errors = append(errors, subErrors...)
				continue
			}
		}
	}

	return errors
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

func NewWrongTypeError(k, e string, v interface{}) error {
	return &WrongTypeError{
		key:          k,
		expectedType: e,
		value:        v,
	}
}
func (w *WrongTypeError) Error() string {
	return fmt.Sprintf("key %v has wrong type %T (should be %s)", w.key, w.value, w.expectedType)
}
