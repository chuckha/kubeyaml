package validation

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestSwagger(t *testing.T) {
	t.Run("a successfully loaded swagger object should", func(t *testing.T) {
		s := loadSwagger(t)
		t.Run("return a valid schema from a version & kind", func(t *testing.T) {
			schema, err := s.FromVersionKind("v1", "Pod")
			if err != nil {
				t.Fatal(err)
			}
			if schema.GVK[0].Kind != "Pod" {
				t.Fatal("did not return a Pod kind")
			}
		})

		t.Run("return a valid schema from a ref", func(t *testing.T) {
			schema, err := s.ForRef("io.k8s.api.events.v1beta1.Event")
			if err != nil {
				t.Fatal(err)
			}
			if schema.GVK[0].Kind != "Event" {
				t.Fatal("did not return EventSeries kind")
			}
		})

		t.Run("return an error for an unknown version & kind", func(t *testing.T) {
			if _, err := s.FromVersionKind("fake", "super fake"); err == nil {
				t.Fatal("there should have been an error but there was not")
			}
		})

		t.Run("validate valid yaml", func(t *testing.T) {
			data := loadData(t, "pod-valid.yaml")
			schema, _ := s.FromVersionKind("v1", "Pod")
			if errs := s.Validate(data, schema, []string{}); len(errs) != 0 {
				t.Fatal("should not have been any errors, but got", errs)
			}
		})

		t.Run("fail to validate invalid yaml", func(t *testing.T) {
			data := loadData(t, "pod-invalid-spec.yaml")
			schema, _ := s.FromVersionKind("v1", "Pod")
			if errs := s.Validate(data, schema, []string{}); len(errs) != 2 {
				t.Fatal("have got 2 errors but got", errs)
			}
		})
	})
}

// Load reads the input and returns the internal type representing the top level document
// that is properly cleaned.
func loadData(t *testing.T, filename string) map[interface{}]interface{} {
	b, err := ioutil.ReadFile("testdata/" + filename)
	if err != nil {
		t.Fatal(err)
	}
	incoming := map[interface{}]interface{}{}
	if err := yaml.Unmarshal(b, incoming); err != nil {
		t.Fatal(err)
	}
	return incoming
}

func loadSwagger(t *testing.T) *Swagger {
	data, err := ioutil.ReadFile("data/swagger-1.19.json")
	if err != nil {
		t.Fatal(err)
	}
	swagger := &Swagger{}
	if err := json.Unmarshal(data, swagger); err != nil {
		t.Fatal(err)
	}
	return swagger
}
