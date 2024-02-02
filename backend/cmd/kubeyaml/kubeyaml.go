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

	"github.com/cristifalcas/kubeyaml/backend/internal/kubernetes"
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
	if err := run(os.Stdin, os.Stdout, os.Args[1:]...); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run(in io.Reader, out io.Writer, args ...string) error {
	opts := &options{}
	validate := flag.NewFlagSet("validate", flag.ExitOnError)
	opts.versions = validate.String("versions", "1.29,1.28,1.27,1.26,1.25,1.24,1.23,1.22,1.21,1.20", "comma separated list of kubernetes versions to validate against")
	opts.silent = validate.Bool("silent", false, "if true, kubeyaml will not print any output")
	validate.Parse(args)
	err := opts.Validate()
	if err != nil {
		return &mainError{"unable to validate options input", err}
	}

	loader := kubernetes.NewLoader()
	gf := kubernetes.NewAPIKeyer("io.k8s.api", ".k8s.io")

	// Read the input
	reader := bufio.NewReader(in)
	var input bytes.Buffer
	readerCopy := io.TeeReader(reader, &input)
	i, err := loader.Load(readerCopy)
	if err != nil {
		return &mainError{input.String(), err}
	}

	aggregatedErrors := &aggErr{}
	for _, version := range opts.Versions {
		reslover, err := kubernetes.NewResolver(version)
		if err != nil {
			aggregatedErrors.Add(fmt.Errorf("%s: %v", version, err))
			continue
		}
		validator := kubernetes.NewValidator(reslover)

		schema, err := reslover.Resolve(gf.APIKey(i.APIVersion, i.Kind))
		if err != nil {
			aggregatedErrors.Add(fmt.Errorf("%s: %v", version, err))
			continue
		}

		if len(aggregatedErrors.errors) > 0 {
			return aggregatedErrors
		}

		errors := validator.Validate(i.Data, schema)
		if len(errors) > 0 {
			if !*opts.silent {
				fmt.Fprintln(out, string(redbg(errors[0].Error())))
				fmt.Fprintln(out, colorize(errors[0], input.Bytes()))
			}
			return &aggErr{errors}
		}
	}
	return nil
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

type mainError struct {
	message string
	err     error
}

func (m *mainError) Error() string {
	return fmt.Sprintf("%s\n%s", m.message, m.err.Error())
}

type aggErr struct {
	errors []error
}

func (a *aggErr) Error() string {
	out := []string{}
	for _, err := range a.errors {
		out = append(out, err.Error())
	}
	return strings.Join(out, "\n")
}
func (a *aggErr) Add(err error) {
	a.errors = append(a.errors, err)
}
