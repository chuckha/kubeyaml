package kubernetes

import "strings"

// APIKeyer can look up a schema key given the apiVersion and the kind.
type APIKeyer struct {
	namespace string
	suffix    string
}

// NewAPIKeyer returns an APIKeyer that finds keys with a common prefix.
func NewAPIKeyer(ns, suffix string) *APIKeyer {
	return &APIKeyer{
		namespace: ns,
		suffix:    suffix,
	}
}

// APIKey returns the key of the API object as listed in the swagger definition.
func (g *APIKeyer) APIKey(apiVersion, kind string) string {
	if apiVersion == "v1" {
		apiVersion = "core/v1"
	}

	// TODO(chuckha) inefficient hack, this is spliiting group/version into
	// group and version then stripping the suffix off the group and rejoining.
	gv := strings.Split(apiVersion, "/")
	if len(gv) == 2 {
		gv[0] = strings.TrimSuffix(gv[0], g.suffix)
		apiVersion = strings.Join(gv, "/")
	}

	apiVersion = strings.Replace(apiVersion, "/", ".", -1)

	return strings.Join([]string{g.namespace, apiVersion, kind}, ".")
}
