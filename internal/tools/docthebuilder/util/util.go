package util

import (
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
