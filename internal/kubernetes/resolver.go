package kubernetes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

// Resolver looks up a schema based on the schema key.
type Resolver struct {
	version string
	swagger *Swagger
}

// NewResolver loads a swagger file, keeps track of the version and returns an instantiated resolver.
func NewResolver(version string) (*Resolver, error) {
	file := fmt.Sprintf("internal/kubernetes/data/swagger-%v.json", version)
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %v: %v", file, err)
	}
	swagger := &Swagger{}
	if err := json.Unmarshal(b, swagger); err != nil {
		return nil, fmt.Errorf("failed to unmarshal swagger file: %v", err)
	}

	return &Resolver{
		version: version,
		swagger: swagger,
	}, nil
}

func (r *Resolver) Resolve(schemaKey string) (*Schema, error) {
	schemaKey = strings.TrimPrefix(schemaKey, "#/definitions/")
	def, ok := r.swagger.Definitions[schemaKey]
	if !ok {
		// Ideally we return the passed in apiVersion but that requires a reshuffling of some functions
		return nil, NewYamlPathError([]string{"apiVersion"}, "", NewUnknownSchemaError(schemaKey))
	}
	return def, nil
}
func (r *Resolver) Version() string {
	return r.version
}
