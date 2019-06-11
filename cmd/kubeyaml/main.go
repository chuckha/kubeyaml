package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/chuckha/kubeyaml/internal/kubernetes"
)

/*
kubeyaml - (read from standard in until EOF, output the same yaml colorized)
output

	filename
contents
..append

--silent => do not output anything just exit 0 on success 1 on failure
--versions 1.8,1.9,1.10,1.11,1.12
*/

type options struct {
	Versions []string
	versions *string
	silent   *bool
}

func (o *options) Validate() error {
	if o.versions == nil {
		return errors.New("versions cannot be nil")
	}
	o.Versions = strings.Split(*o.versions, ",")
	return nil
}

func main() {
	opts := &options{}
	validate := flag.NewFlagSet("validate", flag.ExitOnError)
	opts.versions = validate.String("versions", "1.14,1.13,1.12,1.11,1.10,1.9,1.8", "comma separated list of kubernetes versions to validate against")
	opts.silent = validate.Bool("silent", false, "if true, kubeyaml will not print any output")
	validate.Parse(os.Args[1:])
	err := opts.Validate()
	if err != nil {
		fmt.Println("unable to validate options input")
		os.Exit(1)
	}

	loader := kubernetes.NewLoader()
	gf := kubernetes.NewAPIKeyer("io.k8s.api", ".k8s.io")

	// Read the input
	reader := bufio.NewReader(os.Stdin)
	var input bytes.Buffer
	readerCopy := io.TeeReader(reader, &input)
	i, err := loader.Load(readerCopy)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	exitCode := 0

	for _, version := range opts.Versions {
		reslover, err := kubernetes.NewResolver(version)
		if err != nil {
			fmt.Printf("%s: %v\n", version, err)
			continue
		}
		validator := kubernetes.NewValidator(reslover)

		schema, err := reslover.Resolve(gf.APIKey(i.APIVersion, i.Kind))
		if err != nil {
			fmt.Printf("%s: %v\n", version, err)
			continue
		}

		errors := validator.Validate(i.Data, schema)
		if len(errors) > 0 {
			if !*opts.silent {
				fmt.Println(string(redbg(errors[0].Error())))
				fmt.Println(colorize(errors[0], input.Bytes()))
			}
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}

func colorize(err error, value []byte) string {
	yamlPathErr, ok := err.(*kubernetes.YamlPathError)
	if !ok {
		return "not a *kubernetes.YamlPathError"
	}
	keys := strings.Split(yamlPathErr.Path, ".")
	val, ok := yamlPathErr.Value.(string)
	if !ok {
		return "not a string value"
	}

	for i, key := range keys[:len(keys)-2] {
		re := regexp.MustCompile(fmt.Sprintf("(%s%s):", strings.Repeat("[ -] ", i), key))
		value = re.ReplaceAll(value, red("$1"))
	}
	re := regexp.MustCompile(fmt.Sprintf(`(%s%s:\s*"?%s\"?)`, strings.Repeat("[ -] ", len(keys)-1), keys[len(keys)-1], val))
	value = re.ReplaceAll(value, red("$1"))
	return string(value)
}

func red(in string) []byte {
	return []byte(fmt.Sprintf("\x1b[31;1m%s\x1b[0m", in))
}
func redbg(in string) []byte {
	return []byte(fmt.Sprintf("\x1b[97;1m\x1b[41;1m%s\x1b[0m", in))
}
