package kubernetes

import "fmt"

type Swagger struct {
	Definitions map[string]*Schema
}

type Schema struct {
	Description  string
	Required     []string
	Properties   map[string]*Property
	Initializers *Initializers
	Kind         *Kind
	Metadata     *Metadata
	GVK          []*GroupVersionKind `json:"x-kubernetes-group-version-kind"`
}

type Property struct {
	APIVersion           string
	Description          string
	Type                 string
	Format               string
	Items                *Items
	AdditionalProperties AdditionalProperties
	Reference            string `json:"$ref"`
}

func (p *Property) String() string {
	return fmt.Sprintf(`
APIVersion: %s
Type: %s
Items: %v
Reference: %s
`, p.APIVersion, p.Type, p.Items, p.Reference)
}

type Initializers struct {
	Description   string
	Type          string
	Items         *Items
	PatchMergeKey string
	PatchStrategy string
}

type Items struct {
	Description string
	Type        string
	Reference   string `json:"$ref"`
	Items       *Items
}

type Kind struct {
	Description string
	Type        string
}

type Metadata struct {
	Description string
	Reference   string `json:"$ref"`
}

type GroupVersionKind struct {
	Group   string
	Kind    string
	Version string
}

type AdditionalProperties struct {
	Type string
}
