package validation

import (
	"encoding/json"
	"fmt"

	"github.com/cristifalcas/kubeyaml/backend/internal/kubernetes/data"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type defaultSwaggerVersions struct {
	versions map[string]*Swagger
}

func newDefaultSwaggerVersions() (*defaultSwaggerVersions, error) {
	versions := make(map[string]*Swagger)
	swaggers, err := data.AllStaticFiles()
	if err != nil {
		return nil, err
	}
	for v, swagger := range swaggers {
		swagg := &Swagger{}
		if err := json.Unmarshal(swagger, swagg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal swagger file: %v", err)
		}
		versions[v] = swagg
	}
	return &defaultSwaggerVersions{
		versions: versions,
	}, nil
}

func (d *defaultSwaggerVersions) For(k8sVersion string) (*Swagger, error) {
	out, ok := d.versions[k8sVersion]
	if !ok {
		return nil, errors.Errorf("unknown kubernetes version requested: %q", k8sVersion)
	}
	return out, nil
}

type yamlLoader struct{}

func newYAMLLoader() *yamlLoader {
	return &yamlLoader{}
}

func (d *yamlLoader) Load(in []byte) (*Input, error) {
	incoming := map[interface{}]interface{}{}
	if err := yaml.Unmarshal(in, incoming); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml with error %v", err)
	}

	val, ok := incoming["apiVersion"]
	if !ok {
		return nil, NewRequiredKeyNotFoundError("apiVersion", []string{"apiVersion"})
	}

	apiVersion, ok := val.(string)
	if !ok {
		return nil, NewYamlPathError([]string{"apiVersion"}, val, NewUnknownTypeError(val))
	}

	val, ok = incoming["kind"]
	if !ok {
		return nil, NewRequiredKeyNotFoundError("kind", []string{"kind"})
	}

	kind, ok := val.(string)
	if !ok {
		return nil, NewYamlPathError([]string{"kind"}, val, NewUnknownTypeError(val))
	}

	delete(incoming, "apiVersion")
	delete(incoming, "kind")
	return &Input{
		APIVersion: apiVersion,
		Kind:       kind,
		Data:       incoming,
	}, nil
}
