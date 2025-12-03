// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package util

import (
	"log"
	"strings"

	"github.com/spf13/afero"
)

func FileExists(fs afero.Fs, path string) bool {
	exists, _ := afero.Exists(fs, path)
	return exists
}

func DirExists(fs afero.Fs, path string) bool {
	exists, _ := afero.DirExists(fs, path)
	return exists
}

func MapKeys2Slice[K comparable, V any](m map[K]V) []K {
	result := make([]K, 0, len(m))
	for k := range m {
		result = append(result, k)
	}

	return result
}

func MapValues2Slice[K comparable, V any](m map[K]V) []V {
	result := make([]V, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}

	return result
}

func Map2Slices[K comparable, V any](m map[K]V) ([]K, []V) {
	keys, values := make([]K, 0, len(m)), make([]V, 0, len(m))
	for k, v := range m {
		keys = append(keys, k)
		values = append(values, v)
	}

	return keys, values
}

func FirstCodeValue(line string) string {
	if vals := extractCodeValue(line); len(vals) > 0 {
		return vals[0]
	}
	return ""
}

func extractCodeValue(line string) (res []string) {
	idx1 := strings.Index(line, "`")
	for idx1 >= 0 {
		idx2 := idx1 + 1 + strings.Index(line[idx1+1:], "`")
		if idx2 > len(line) || idx2 <= idx1 {
			log.Printf("ExtractCodeValue: code mark ` not closed in '%s'", line)
			return
		}
		res = append(res, line[idx1+1:idx2])
		nextIdx := strings.Index(line[idx2+1:], "`")
		if nextIdx < 0 {
			break
		}
		idx1 = idx2 + 1 + nextIdx
	}
	return
}
