package internal

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// only deal with major & minor
type version struct {
	major, minor int
}

func (v version) String() string {
	return fmt.Sprintf("%d.%d", v.major, v.minor)
}

type versions []version

func (v versions) Len() int { return len(v) }
func (v versions) Less(i, j int) bool {
	if v[i].major > v[j].major {
		return true
	}
	if v[i].major < v[j].major {
		return false
	}
	// at this point majors are the same
	if v[i].minor > v[j].minor {
		return true
	}
	return false
}
func (v versions) Swap(i, j int) { v[i], v[j] = v[j], v[i] }

// SortVersions sorts a list of versions from newest to oldest
func SortVersions(vs ...string) ([]string, error) {
	internalVersions := make(versions, len(vs))
	out := make([]string, len(vs))
	for i, v := range vs {
		parts := strings.Split(v, ".")
		if len(parts) != 2 {
			return out, errors.Errorf("Invalid version: %q, requires format x.y", v)
		}
		major, err := strconv.Atoi(parts[0])
		if err != nil {
			return out, errors.WithMessagef(err, "major part of version is not an int: %q", parts[0])
		}
		minor, err := strconv.Atoi(parts[1])
		if err != nil {
			return out, errors.WithMessagef(err, "minor part of version is not an int: %q", parts[0])
		}
		internalVersions[i] = version{major, minor}
	}
	sort.Sort(internalVersions)
	for i, v := range internalVersions {
		out[i] = v.String()
	}
	return out, nil
}
