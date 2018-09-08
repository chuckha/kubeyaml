package kubernetes

type Validator struct{}

func (v *Validator) Validate(incoming map[string]interface{}, schema *Schema) error {
	return nil
}
