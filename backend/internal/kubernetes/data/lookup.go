package data

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

func AllStaticFiles() (map[string][]byte, error) {
	out := make(map[string][]byte)
	infos, err := ioutil.ReadDir(".")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for _, fi := range infos {
		if strings.HasSuffix(fi.Name(), ".json") {
			data, err := ioutil.ReadFile(fi.Name())
			if err != nil {
				return nil, errors.Wrapf(err, "filename: %s", fi.Name())
			}
			version := strings.TrimPrefix(fi.Name(), "swagger-")
			version = strings.TrimSuffix(version, ".json")
			out[version] = data
		}
	}
	return out, nil
}

type StaticFiles struct{}

// Swagger is a fairly meh function. It's poorly named and tied to the update-schemas file.
func (s *StaticFiles) Swagger(version string) []byte {
	switch version {
	case "1.20":
		return s.OneTwenty()
	case "1.21":
		return s.OneTwentyone()
	case "1.22":
		return s.OneTwentytwo()
	case "1.23":
		return s.OneTwentythree()
	case "1.24":
		return s.OneTwentyfour()
	case "1.25":
		return s.OneTwentyfive()
	case "1.26":
		return s.OneTwentysix()
	case "1.27":
		return s.OneTwentyseven()
	case "1.28":
		return s.OneTwentyeight()
	case "1.29":
		return s.OneTwentynine()
	default:
		panic(fmt.Sprintf("unknown version %v", version))
	}
}
