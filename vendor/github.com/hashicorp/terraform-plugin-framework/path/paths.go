// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package path

import "strings"

// Paths is a collection of exact attribute paths.
//
// Refer to the Path documentation for more details about intended usage.
type Paths []Path

// Append adds the given Paths to the collection without duplication and
// returns the combined result.
func (p *Paths) Append(paths ...Path) Paths {
	if p == nil {
		return paths
	}

	for _, newPath := range paths {
		if p.Contains(newPath) {
			continue
		}

		*p = append(*p, newPath)
	}

	return *p
}

// Contains returns true if the collection of paths includes the given path.
func (p Paths) Contains(checkPath Path) bool {
	for _, path := range p {
		if path.Equal(checkPath) {
			return true
		}
	}

	return false
}

// String returns the human-readable representation of the path collection.
// It is intended for logging and error messages and is not protected by
// compatibility guarantees.
//
// Empty paths are skipped.
func (p Paths) String() string {
	var result strings.Builder

	result.WriteString("[")

	for pathIndex, path := range p {
		if path.Equal(Empty()) {
			continue
		}

		if pathIndex != 0 {
			result.WriteString(",")
		}

		result.WriteString(path.String())
	}

	result.WriteString("]")

	return result.String()
}
