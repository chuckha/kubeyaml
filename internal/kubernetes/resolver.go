package kubernetes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type Resolver struct {
	swagger *Swagger
}

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
		swagger: swagger,
	}, nil
}

type UnknownSchemaError struct {
	schemaKey string
}

func NewUnknownSchemaError(key string) error {
	return &UnknownSchemaError{schemaKey: key}
}
func (u *UnknownSchemaError) Error() string {
	return fmt.Sprintf("Unknown schema %v", u.schemaKey)
}

func (r *Resolver) Resolve(schemaKey string) (*Schema, error) {
	schemaKey = strings.TrimPrefix(schemaKey, "#/definitions/")
	def, ok := r.swagger.Definitions[schemaKey]
	if !ok {
		return nil, NewUnknownSchemaError(schemaKey)
	}
	return def, nil
}
