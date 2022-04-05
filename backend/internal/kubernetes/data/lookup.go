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
	case "1.8":
		return s.OneEight()
	case "1.9":
		return s.OneNine()
	case "1.10":
		return s.OneTen()
	case "1.11":
		return s.OneEleven()
	case "1.12":
		return s.OneTwelve()
	case "1.13":
		return s.OneThirteen()
	case "1.14":
		return s.OneFourteen()
	case "1.15":
		return s.OneFifteen()
	case "1.16":
		return s.OneSixteen()
	case "1.17":
		return s.OneSeventeen()
	case "1.18":
		return s.OneEighteen()
	case "1.19":
		return s.OneNineteen()
	case "1.20":
		return s.OneTwenty()
	case "1.21":
		return s.OneTwentyone()
	case "1.22":
		return s.OneTwentytwo()
	default:
		panic(fmt.Sprintf("unknown version %v", version))
	}
}
