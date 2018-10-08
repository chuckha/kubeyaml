package kubernetes

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chuckha/kubeyaml/internal/kubernetes/data"
)

// Resolver looks up a schema based on the schema key.
type Resolver struct {
	version string
	swagger *Swagger
}

// NewResolver loads a swagger file, keeps track of the version and returns an instantiated resolver.
func NewResolver(version string) (*Resolver, error) {
	staticFiles := &data.StaticFiles{}
	swagger := &Swagger{}
	if err := json.Unmarshal(staticFiles.Swagger(version), swagger); err != nil {
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
