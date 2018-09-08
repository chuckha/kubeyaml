package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Run this from the top level dir
func main() {
	releases := []string{
		"1.8",
		"1.9",
		"1.10",
		"1.11",
		"1.12",
	}
	urlFmt := "https://raw.githubusercontent.com/kubernetes/kubernetes/release-%s/api/openapi-spec/swagger.json"

	for _, release := range releases {
		r, err := http.Get(fmt.Sprintf(urlFmt, release))
		if err != nil {
			fmt.Printf("failed to get schemas for release %v: %v\n", release, err)
			continue
		}

		schema, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("headers", r.Header)
			fmt.Println("status code:", r.StatusCode)
			fmt.Printf("failed to read response body: %v\n", err)
			continue
		}

		outDir, err := os.Getwd()
		if err != nil {
			fmt.Printf("could not get working dir: %v", err)
			continue
		}

		if err := ioutil.WriteFile(fmt.Sprintf("%s/internal/kubernetes/data/swagger-%s.json", outDir, release), schema, os.FileMode(uint32(0644))); err != nil {
			fmt.Printf("error writing file release-%s: %v", release, err)
			continue
		}
	}
}
