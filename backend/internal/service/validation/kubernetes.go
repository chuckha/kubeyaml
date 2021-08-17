package validation

import (
	"fmt"
	"strings"
	"time"
)

// Swagger is a map of schema keys to schema objects loaded in from a swagger file.
type Swagger struct {
	Definitions map[string]*Schema
}

func (s *Swagger) ForRef(ref string) (*Schema, error) {
	ref = strings.TrimPrefix(ref, "#/definitions/")
	def, ok := s.Definitions[ref]
	if !ok {
		// Ideally we return the passed in apiVersion but that requires a reshuffling of some functions
		return nil, NewYamlPathError([]string{"apiVersion"}, "", NewUnknownSchemaError(ref))
	}
	return def, nil
}

// APIKey returns the key of the API object as listed in the swagger definition.
func (s *Swagger) FromVersionKind(apiVersion, kind string) (*Schema, error) {
	namespace := "io.k8s.api"
	suffix := ".k8s.io"
	if apiVersion == "v1" {
		apiVersion = "core/v1"
	}
	gv := strings.Split(apiVersion, "/")
	if len(gv) == 2 {
		gv[0] = strings.TrimSuffix(gv[0], suffix)
		apiVersion = strings.Join(gv, "/")
	}

	apiVersion = strings.Replace(apiVersion, "/", ".", -1)

	return s.ForRef(strings.Join([]string{namespace, apiVersion, kind}, "."))
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
func (s *Swagger) Validate(incoming map[interface{}]interface{}, schema *Schema, path []string) []error {
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
				schema, err := s.ForRef(property.Items.Reference)
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

					if errs := s.handleObject(key, item, tlpIdx, schema); len(errs) > 0 {
						errors = append(errors, errs...)
					}
				}
			}
		// default is some k8s object
		default:
			schema, err := s.ForRef(property.Reference)
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
			if errs := s.handleObject(key, value, tlp, schema); len(errs) > 0 {
				errors = append(errors, errs...)
			}
		}
	}

	return errors
}

// handleObject takes a key that has a value that will be of type map[interface{}]interface{}
// handleObject takes the current path to the key that is being validated and the schema of the object hidden under the value interface.
func (s *Swagger) handleObject(key string, value interface{}, path []string, schema *Schema) []error {
	object, ok := value.(map[interface{}]interface{})
	if !ok {
		return []error{NewYamlPathError(path, value, NewWrongTypeError(key, "map[interface{}]interface{}", value))}
	}
	return s.Validate(object, schema, path)
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

// Schema is a swagger schema. I'm sure there's a real definition somewhere but this gets everything this program needs.
type Schema struct {
	// Description is the description of the object that conforms to this schema.
	Description string
	// Required lists the required fields of this object.
	Required []string
	// Properties
	Properties map[string]*Property
	// Initializers are optional. TODO(chuckha) what are these actually?
	Initializers *Initializers
	// Kind is an optional field that describes the kind of object inside a list type
	Kind *Kind
	// Metadata is an optional field that contains a reference to some other schema.
	Metadata *Metadata
	// GVK is the GroupVersionKind from kubernetes.
	GVK []*GroupVersionKind `json:"x-kubernetes-group-version-kind"`
	// Type is when the object is actually a type rename  of a builtin (type X string)
	Type string
	// Format is the format of the type when Type is "string"
	Format string
}

// Property is a single property, or field, of a schema.
type Property struct {
	// Description is the description of the property being defined.
	Description string
	// Type is a basic type like string, object, integer or array.
	Type string
	// Format is a sub type, such as string could be `byte` or integer could be `int64`.
	Format string
	// Items is required for array types. This tells us what types are inside the array.
	Items *Items
	// AdditionalProperties is an optional field for object types defining what kind of key/value pairs will be found on the object.
	// TODO(chuckha) make this a pointer type
	AdditionalProperties AdditionalProperties
	// Reference is a reference to another schema in our list of definitions.
	Reference string `json:"$ref"`
}

// String implements the Stringer interface and gives us a nice human readable output.
func (p *Property) String() string {
	return fmt.Sprintf(`
Type: %s
Items: %v
Reference: %s
`, p.Type, p.Items, p.Reference)
}

// Initializers are things that will run before a pod can start running.
// This is only used in two places, the object metadata and the initializersConfiguration object.
// See: https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#how-are-initializers-triggered
type Initializers struct {
	// Description is a description of the initializers.
	Description string
	// Type is the type of initializers, this is always an array.
	Type string
	// Items are the actual list of initializers.
	Items *Items
	// TODO(chuckha) Can probably remove these two
	PatchMergeKey string `json:"x-kubernetes-patch-merge-key"`
	PatchStrategy string `json:"x-kubernetes-patch-strategy"`
}

// Items are the array of items when the type is "array".
type Items struct {
	// Description is a description of the array of items.
	Description string
	// Type is the type of item in the array.
	Type string
	// Reference is the reference to the schema type of items stored in this array.
	Reference string `json:"$ref"`
	// Items can be an array of arrays of arrays of ...
	Items *Items
}

// Kind is the kind we all know and love from ObjectMeta.
type Kind struct {
	// Description describes what kind is, it's an resource in the API.
	Description string
	// Type is always a string in the Kind object.
	Type string
}

// Metadata is a reference to the shared metadata schema.
type Metadata struct {
	// Description descries what metadata is.
	Description string
	// Reference is the key to the metadata schema.
	Reference string `json:"$ref"`
}

// GroupVersionKind is the gvk of kubernetes schema.
type GroupVersionKind struct {
	// Group is the group such as v1, apps, etc.
	Group string
	// Kind is the type of object such as Deployment or Pod.
	Kind string
	// Version is the API version such as v1, v1alpha3, v1beta2.
	Version string
}

// AdditionalProperties define the type found in objects.
// Kubernetes only uses string:string dicts, but this might break on more advanced use cases.
type AdditionalProperties struct {
	// Type will almost always be a string
	Type string
}
