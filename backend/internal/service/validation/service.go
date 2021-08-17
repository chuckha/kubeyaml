package validation

type SwaggerService interface {
	Validate(incoming map[interface{}]interface{}, schema *Schema, path []string) []error
	FromVersionKind(apiVersion, kind string) (*Schema, error)
	ForRef(ref string) (*Schema, error)
}

type loader interface {
	Load(data []byte) (*Input, error)
}

// Input is the top level YAML document the system receives.
type Input struct {
	Kind       string
	APIVersion string
	Data       map[interface{}]interface{}
}

type Service struct {
	swagger SwaggerService
	loader  loader
}

type ServiceOption func(s *Service)

func WithSwaggerService(ss SwaggerService) ServiceOption {
	return func(s *Service) {
		s.swagger = ss
	}
}

func WithLoader(l loader) ServiceOption {
	return func(s *Service) {
		s.loader = l
	}
}

func NewService(opts ...ServiceOption) *Service {
	s := &Service{
		loader: newYAMLLoader(),
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

// input is a kubernetes object as yaml
// validate it against a specific type schema
func (s *Service) Validate(input []byte) error {
	loaded, err := s.loader.Load(input)
	if err != nil {
		return err
	}

	// look up the schema for the version
	schema, err := s.swagger.FromVersionKind(loaded.APIVersion, loaded.Kind)
	if err != nil {
		return err
	}
	errs := s.swagger.Validate(loaded.Data, schema, []string{})
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
