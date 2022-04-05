package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const templateFormat = `package data

func (s *StaticFiles) %s() []byte {
	return []byte(` + "`%s`" + `)
}
`

// Run this from the top level dir
func main() {
	releases := map[string]string{
		// "OneEight":    "1.8",
		// "OneNine":     "1.9",
		// "OneTen":      "1.10",
		// "OneEleven":   "1.11",
		// "OneTwelve":   "1.12",
		"OneThirteen":  "1.13",
		"OneFourteen":  "1.14",
		"OneFifteen":   "1.15",
		"OneSixteen":   "1.16",
		"OneSeventeen": "1.17",
		"OneEighteen":  "1.18",
		"OneNineteen":  "1.19",
		"OneTwenty":    "1.20",
		"OneTwentyone": "1.21",
		"OneTwentytwo": "1.22",
	}
	urlFmt := "https://raw.githubusercontent.com/kubernetes/kubernetes/release-%s/api/openapi-spec/swagger.json"

	for funcName, release := range releases {
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
		schema = bytes.Replace(schema, []byte("`"), []byte("` + \"`\" + `"), -1)
		fileData := fmt.Sprintf(templateFormat, funcName, schema)

		if err := ioutil.WriteFile(fmt.Sprintf("%s/internal/kubernetes/data/swagger_%s.go", outDir, funcName), []byte(fileData), os.FileMode(uint32(0644))); err != nil {
			fmt.Printf("error writing go file for release-%s: %v", release, err)
			continue
		}

	}
}
