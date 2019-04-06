package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/chuckha/kubeyaml/internal/kubernetes"
)

func main() {
	f, err := os.Create("dist/index.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	t, err := template.ParseFiles("templates/index.template.html")
	if err != nil {
		panic(err)
	}

	versions := []string{"1.8", "1.9", "1.10", "1.11", "1.12", "1.13"}
	validators := make([]validator, len(versions))
	for i, version := range versions {
		resolver, err := kubernetes.NewResolver(version)
		if err != nil {
			fmt.Printf("failed to get a resolver for version %q: %v", version, err)
			os.Exit(1)
		}
		validators[i] = kubernetes.NewValidator(resolver)
	}

	if err := t.Execute(f, indexTemplateData{
		Validators: validators,
		Selected:   "1.8",
	}); err != nil {
		panic(err)
	}
	fmt.Printf("Generated file: %s\n", "dist/index.html")
}

type validator interface {
	Validate(map[interface{}]interface{}, *kubernetes.Schema) []error
	Resolve(string) (*kubernetes.Schema, error)
	Version() string
}

type indexTemplateData struct {
	Validators []validator
	Selected   string
}
