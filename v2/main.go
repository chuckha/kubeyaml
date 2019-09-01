package main

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/encoding/yaml"
)

func main() {
	var r cue.Runtime
	c, err := yaml.Decode(&r, "testdata/pod.yaml", "test")
	if err != nil {
		panic(err)
	}
	fmt.Println(c)
}
