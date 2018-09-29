package kubernetes_test

import (
	"testing"

	"github.com/chuckha/kube-validate/internal/kubernetes"
)

type resolver struct {
	lookup  *kubernetes.Swagger
	version string
}

func (r *resolver) Resolve(schemaKey string) (*kubernetes.Schema, error) {
	return r.lookup.Definitions[schemaKey], nil
}
func (r *resolver) Version() string {
	return r.version
}

func TestValidate(t *testing.T) {
	// shared resolver
	testcases := []struct {
		name        string
		definitions map[string]*kubernetes.Schema
		version     string
		incoming    map[interface{}]interface{}
		schema      *kubernetes.Schema
		path        []string
		check       func([]error, *testing.T)
	}{
		{
			name: "easy test",
			definitions: map[string]*kubernetes.Schema{
				"webpage": &kubernetes.Schema{},
				"Buttons": &kubernetes.Schema{
					Properties: map[string]*kubernetes.Property{
						"submit": &kubernetes.Property{
							Type: "string",
						},
					},
				},
				"headingNumber": &kubernetes.Schema{
					Properties: map[string]*kubernetes.Property{
						"size": &kubernetes.Property{
							Type: "integer",
						},
					},
				},
			},
			version: "v1.12",
			incoming: map[interface{}]interface{}{
				interface{}("buttons"): interface{}(map[interface{}]interface{}{
					interface{}("submit"): interface{}("blue"),
				}),
				interface{}("headers"): interface{}([]interface{}{
					interface{}(map[interface{}]interface{}{
						interface{}("size"): interface{}(2),
					}),
				}),
			},
			schema: &kubernetes.Schema{
				Properties: map[string]*kubernetes.Property{
					"buttons": &kubernetes.Property{
						Reference: "Buttons",
					},
					"headers": &kubernetes.Property{
						Type: "array",
						Items: &kubernetes.Items{
							Reference: "headingNumber",
						},
					},
				},
			},
			check: func(errors []error, t *testing.T) {
				if len(errors) != 0 {
					t.Fatalf("expected no errors but got %v", errors)
				}
			},
		},
		{
			name: "testing errors",
			definitions: map[string]*kubernetes.Schema{
				"banana": &kubernetes.Schema{},
			},
			version: "v1.1",
			incoming: map[interface{}]interface{}{
				interface{}("sides"):  interface{}("hello"),
				interface{}("other"):  interface{}(4),
				interface{}("yellow"): interface{}(300),
				interface{}("green"):  interface{}("hello"),
				interface{}("red"):    interface{}(23),
				interface{}("blue"):   interface{}("hi"),
			},
			schema: &kubernetes.Schema{
				Properties: map[string]*kubernetes.Property{
					"sides": &kubernetes.Property{
						Type: "integer",
					},
					"other": &kubernetes.Property{
						Type: "string",
					},
					"yellow": &kubernetes.Property{
						Type: "boolean",
					},
					"green": &kubernetes.Property{
						Type: "object",
					},
					"red": &kubernetes.Property{
						Type: "array",
					},
					"blue": &kubernetes.Property{
						Type:      "some obj",
						Reference: "banana",
					},
				},
			},
			check: func(errors []error, t *testing.T) {
				if len(errors) == 0 {
					t.Fatalf("expected lots of errors but found none")
				}
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			res := &resolver{
				lookup: &kubernetes.Swagger{
					Definitions: tc.definitions,
				},
				version: tc.version,
			}
			v := kubernetes.NewValidator(res)
			errors := v.Validate(tc.incoming, tc.schema)
			tc.check(errors, t)
		})
	}
}
