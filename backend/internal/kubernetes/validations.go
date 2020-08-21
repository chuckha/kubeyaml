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

// Validate sees incoming data and validates it against the known schemas.
// This is recursive so it does a depth first search of all key/values.
// TODO(chuckha) turn this into a stack-based dfs search.
func (v *Validator) Validate(incoming map[interface{}]interface{}, schema *Schema) []error {
	return v.validate(incoming, schema, []string{})
}

// validate is the meat and potatoes of this entire application.
// incoming is a list of key/values pairs from a YAML document.
// schema is the schema we expect incoming to validate against
// path is the list of keys we have traversed to get to this object (as this object could be anywhere in a YAML document)
// This function loops through each key doing the following:
// 1. Checking that the key is a string
// 2. Checking that the key is an expected key
// 3. Checking the value is the expected type
// 4. If an object is encountered then the function is recursive.
func (v *Validator) validate(incoming map[interface{}]interface{}, schema *Schema, path []string) []error {
	errors := make([]error, 0)

	// Validate each key one at a time descending as deep and as wide as the key goes.
	for k, value := range incoming {
		// Keep track of where we are
		tlp := make([]string, len(path))
		copy(tlp, path)

		key, ok := k.(string)
		if !ok {
			errors = append(errors, NewYamlPathError(tlp, "", NewKeyNotStringError(k)))
			continue
		}

		// the key is a string so we can now act on it.
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

				for i, item := range items {
					// Sequences keys will have the index following the key
					// e.g. spec.template.containers.1 (the problem is in the second one)
					tlpIdx := append(tlp, fmt.Sprintf("%d", i))
					if len(schema.Required) > 0 {
						if errs := resolveRequiredFields(key, item, schema, tlpIdx); len(errs) > 0 {
							errors = append(errors, errs...)
						}
					}

					if errs := v.handleObject(key, item, tlpIdx, schema); len(errs) > 0 {
						errors = append(errors, errs...)
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

			if len(schema.Required) > 0 {
				if errs := resolveRequiredFields(key, value, schema, tlp); len(errs) > 0 {
					errors = append(errors, errs...)
				}
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
			if errs := v.handleObject(key, value, tlp, schema); len(errs) > 0 {
				errors = append(errors, errs...)
			}
		}
	}

	return errors
}

// handleObject takes a key that has a value that will be of type map[interface{}]interface{}
// handleObject takes the current path to the key that is being validated and the schema of the object hidden under the value interface.
func (v *Validator) handleObject(key string, value interface{}, path []string, schema *Schema) []error {
	object, ok := value.(map[interface{}]interface{})
	if !ok {
		return []error{NewYamlPathError(path, value, NewWrongTypeError(key, "map[interface{}]interface{}", value))}
	}
	return v.validate(object, schema, path)
}

// objectValue will be a map[interface{}]interface{}
func resolveRequiredFields(key string, objectValue interface{}, schema *Schema, path []string) []error {
	mapvalues, ok := objectValue.(map[interface{}]interface{})
	if !ok {
		return []error{NewWrongTypeError(key, "map[interface{}]interface{}", objectValue)}
	}
	requiredKeys := map[string]bool{}
	for _, k := range schema.Required {
		requiredKeys[k] = false
	}
	for k := range mapvalues {
		realKey, ok := k.(string)
		if !ok {
			return []error{NewKeyNotStringError(key)}
		}
		if _, ok := requiredKeys[realKey]; ok {
			requiredKeys[realKey] = true
		}
	}
	for k, found := range requiredKeys {
		if !found {
			path = append(path, key)
			return []error{NewRequiredKeyNotFoundError(k, path)}
		}
	}
	return []error{}
}
