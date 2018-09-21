package kubernetes

import "fmt"

// Swagger is a map of schema keys to schema objects loaded in from a swagger file.
type Swagger struct {
	Definitions map[string]*Schema
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
