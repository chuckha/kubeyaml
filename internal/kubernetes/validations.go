package kubernetes

import (
	"fmt"
	"time"
)

type resolver interface {
	Resolve(string) (*Schema, error)
	Version() string
}

// Validator knows enough to be able to validate a YAML document.
type Validator struct {
	resolver resolver
}

// NewValidator returns an instantiated validator.
func NewValidator(resolver resolver) *Validator {
	return &Validator{
		resolver: resolver,
	}
}

// Resolve wraps the internal resolver's resolve method.
func (v *Validator) Resolve(schemaKey string) (*Schema, error) {
	return v.resolver.Resolve(schemaKey)
}

// Version wraps the internal resolver's version method.
func (v *Validator) Version() string {
	return v.resolver.Version()
}

// Validate is the meat of this code. It sees incoming data and validates it against the known schemas.
// This is recursive so it does a depth first search of all key/values.
// TODO(chuckha) turn this into a stack-based dfs search.
func (v *Validator) Validate(incoming map[string]interface{}, schema *Schema, path []string) []error {
	errors := make([]error, 0)

	// Validate each key one at a time descending as deep and as wide as the key goes.
	for key, value := range incoming {
		// Keep track of where we are
		tlp := make([]string, len(path))
		copy(tlp, path)
		tlp = append(tlp, key)

		property, ok := schema.Properties[key]
		if !ok {
			errors = append(errors, NewYamlPathError(tlp, "", NewUnknownKeyError(key)))
			continue
		}

		switch property.Type {
		case "string":
			// TODO: formats?
			if _, ok := value.(string); !ok {
				errors = append(errors, NewYamlPathError(tlp, value, NewWrongTypeError(key, "string", value)))
			}
		case "integer":
			// ignore property.Format until it causes a bug
			if _, ok := value.(int); !ok {
				errors = append(errors, NewYamlPathError(tlp, value, NewWrongTypeError(key, "integer", value)))
			}
		case "boolean":
			if _, ok := value.(bool); !ok {
				errors = append(errors, NewYamlPathError(tlp, value, NewWrongTypeError(key, "boolean", value)))
			}
		case "object":
			// this is for things like labels; map[interface{}]interface{} looks weird but that's how our yaml parser works.
			if _, ok := value.(map[interface{}]interface{}); !ok {
				errors = append(errors, NewYamlPathError(tlp, value, NewWrongTypeError(key, "map[interface{}]interface{}", value)))
			}
		case "array":
			items, ok := value.([]interface{})
			if !ok {
				errors = append(errors, NewYamlPathError(tlp, value, NewWrongTypeError(key, "[]interface{}", value)))
				continue
			}

			switch property.Items.Type {
			case "string":
				for _, item := range items {
					if _, ok := item.(string); !ok {
						errors = append(errors, NewWrongTypeError(key, "string", item))
					}
				}
			// assume it's an array of objects
			default:
				// TODO: check that items is not nil
				schema, err := v.resolver.Resolve(property.Items.Reference)
				if err != nil {
					fmt.Println(key, property)
					errors = append(errors, NewYamlPathError(tlp, schema, err))
					continue
				}

				for _, item := range items {
					m, ok := item.(map[interface{}]interface{})
					if !ok {
						errors = append(errors, NewYamlPathError(tlp, item, NewWrongTypeError(key, "map[interface{}]interface{}", item)))
						continue
					}
					converted, err := keysToStrings(m)
					if err != nil {
						errors = append(errors, NewYamlPathError(tlp, item, err))
						continue
					}
					if errs := v.Validate(converted, schema, tlp); len(errs) > 0 {
						errors = append(errors, errs...)
						continue
					}
				}
			}
		// default is some k8s object
		default:
			schema, err := v.resolver.Resolve(property.Reference)
			if err != nil {
				// DEBUG LINE good to use if there is a weird error
				// fmt.Println(key, property)
				errors = append(errors, NewYamlPathError(tlp, property.Reference, err))
				continue
			}

			// Bail if the object reference is a type rename
			if schema.Type == "string" {
				// format must be set if type is string
				switch schema.Format {
				case "int-or-string":
					if _, ok := value.(string); !ok {
						if _, ok2 := value.(int); !ok2 {
							errors = append(errors, NewYamlPathError(tlp, value, NewWrongTypeError(key, "int-or-string", value)))
						}
					}
				case "date-time":
					// nil is valid object reference
					if value == nil {
						continue
					}
					date, ok := value.(string)
					if !ok {
						errors = append(errors, NewYamlPathError(tlp, value, NewWrongTypeError(key, "string", value)))
						continue
					}
					if _, err := time.Parse("2006-01-02T15:04:05Z", date); err != nil {
						errors = append(errors, NewYamlPathError(tlp, value, NewWrongTypeError(key, "time.Time", value)))
					}
				default:
					errors = append(errors, NewYamlPathError(tlp, value, NewUnknownFormatError(schema.Format)))
				}
				continue
			}

			// It's an object if it's not a string.
			d, ok := value.(map[interface{}]interface{})
			if !ok {
				errors = append(errors, NewYamlPathError(tlp, value, NewWrongTypeError(key, "map[interface{}]interface{}", value)))
				continue
			}
			convertedMap, err := keysToStrings(d)
			if err != nil {
				errors = append(errors, NewYamlPathError(tlp, value, err))
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
