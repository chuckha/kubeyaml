package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/chuckha/kubeyaml/internal/kubernetes"
)

func main() {
	loader := kubernetes.NewLoader()
	gf := kubernetes.NewAPIKeyer("io.k8s.api", ".k8s.io")
	versions := []string{"1.8", "1.9", "1.10", "1.11", "1.12"}

	// Read the input
	reader := bufio.NewReader(os.Stdin)
	i, err := loader.Load(reader)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, version := range versions {
		fmt.Println(version)
		reslover, err := kubernetes.NewResolver(version)
		if err != nil {
			fmt.Println(err)
			continue
		}
		validator := kubernetes.NewValidator(reslover)

		schema, err := reslover.Resolve(gf.APIKey(i.APIVersion, i.Kind))
		if err != nil {
			fmt.Println(err)
			continue
		}

		errors := validator.Validate(i.Data, schema)
		fmt.Println(errors)
	}
}
