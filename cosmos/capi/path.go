package capi

import (
	"fmt"
	"strings"
)

type PathBase struct {
	Path  string
	Parts map[string]string
}

func ParsePath(path string) (PathBase, error) {
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	parts := strings.Split(path, "/")

	if len(parts)%2 != 0 {
		return PathBase{}, fmt.Errorf("Unable to parse cosmos path, odd number of parts")
	}

	m := map[string]string{}
	for i := 0; i < len(parts); i += 2 {
		k := parts[i]
		v := parts[i+1]

		// Check key/value for empty strings.
		if k == "" || v == "" {
			return PathBase{}, fmt.Errorf("Key/Value cannot be empty. Key: '%s', Value: '%s'", k, v)
		}

		m[k] = v
	}

	return PathBase{
		Path:  path,
		Parts: m,
	}, nil
}

type IPath interface {
	GetPath() string
	GetCreatePath() string
}

func (p PathBase) GetPath() string {
	return p.Path
}

//given a path (dbs/data/colls/stuff) return the creation path (dbs/data/colls)
func (p PathBase) GetCreatePath() string {
	return p.Path[0:strings.LastIndex(p.Path, "/")]
}
