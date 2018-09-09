package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/chuckha/k8svalidator/internal/kubernetes"
)

func main() {
	loader := kubernetes.NewLoader()
	gf := kubernetes.NewKubernetesGroupFinder("io.k8s.api")
	reslover, err := kubernetes.NewResolver("1.12")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	validator := kubernetes.NewValidator(reslover)

	// Read the input
	reader := bufio.NewReader(os.Stdin)
	i, err := loader.Load(reader)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	schema, err := reslover.Resolve(gf.GroupFind(i.APIVersion, i.Kind))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	errors := validator.Validate(i.Data, schema)
	fmt.Println(errors)
}
