package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/chuckha/k8svalidator/internal/kubernetes"
)

func main() {
	loader := kubernetes.NewLoader()
	//	gf := kubernetes.NewKubernetesGroupFinder("io.k8s.api")
	_, err := kubernetes.NewResolver("1.12")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Read the input
	reader := bufio.NewReader(os.Stdin)
	i, err := loader.Load(reader)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(i)
}

// schema resolver
// given a string return the schema
