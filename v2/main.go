package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"cuelang.org/go/cue"
	"cuelang.org/go/encoding/gocode"
)

func main() {
	// filepath.Walk("pkg", func(path string, info os.FileInfo, err error) error {
	// 	if strings.HasSuffix("types_go_gen.cue")
	// })

	f, _ := os.Open("./pkg/k8s.io/api/core/v1/types_go_gen.cue")
	var r cue.Runtime

	instance, err := r.Compile("testdata/pod.cue", f)
	if err != nil {
		fmt.Println("failed to compile", err)
		return
	}

	b, err := gocode.Generate("k8s.io/api/core/v1", instance, nil)
	if err != nil {
		fmt.Println("failed to generate", err)
		return
	}
	if err := ioutil.WriteFile("cue_gen.go", b, 0644); err != nil {
		fmt.Println("failed to write file", err)
		return
	}

	// c, err := yaml.Decode(&r, "testdata/pod.yaml", "test")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(c)
}
