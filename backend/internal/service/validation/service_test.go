package validation

import (
	"testing"
)

func TestValidation(t *testing.T) {
	t.Run("should return no errors if no errors are found validating", func(t *testing.T) {
		svc := NewService(
			WithSwaggerService(&noErrorDummy{}),
			WithLoader(&dummyLoader{}),
		)
		if err := svc.Validate([]byte{}); err != nil {
			t.Fatal(err)
		}
	})
}

type noErrorDummy struct{}

func (n *noErrorDummy) Validate(incoming map[interface{}]interface{}, schema *Schema, path []string) []error {
	return nil
}
func (n *noErrorDummy) FromVersionKind(apiVersion, kind string) (*Schema, error) { return nil, nil }
func (n *noErrorDummy) ForRef(ref string) (*Schema, error)                       { return nil, nil }

type dummyLoader struct{}

func (d *dummyLoader) Load([]byte) (*Input, error) {
	return &Input{}, nil
}
