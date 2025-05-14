// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
	"maps"
	"slices"
	"sort"
)

var _ Check = objectExact{}

type objectExact struct {
	value map[string]Check
}

// CheckValue determines whether the passed value is of type map[string]any, and
// contains matching object entries.
func (v objectExact) CheckValue(other any) error {
	otherVal, ok := other.(map[string]any)

	if !ok {
		return fmt.Errorf("expected map[string]any value for ObjectExact check, got: %T", other)
	}

	if len(otherVal) != len(v.value) {
		deltaMsg := ""
		if len(otherVal) > len(v.value) {
			deltaMsg = createDeltaString(otherVal, v.value, "actual value has extra attribute(s): ")
		} else {
			deltaMsg = createDeltaString(v.value, otherVal, "actual value is missing attribute(s): ")
		}

		return fmt.Errorf("expected %d attribute(s) for ObjectExact check, got %d attribute(s): %s", len(v.value), len(otherVal), deltaMsg)
	}

	var keys []string

	for k := range v.value {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, k := range keys {
		otherValItem, ok := otherVal[k]

		if !ok {
			return fmt.Errorf("missing attribute %s for ObjectExact check", k)
		}

		if err := v.value[k].CheckValue(otherValItem); err != nil {
			return fmt.Errorf("%s object attribute: %s", k, err)
		}
	}

	return nil
}

// String returns the string representation of the value.
func (v objectExact) String() string {
	var keys []string

	for k := range v.value {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	mapVals := make(map[string]string, len(keys))

	for _, k := range keys {
		mapVals[k] = v.value[k].String()
	}

	return fmt.Sprintf("%v", mapVals)
}

// ObjectExact returns a Check for asserting equality between the supplied
// map[string]Check and the value passed to the CheckValue method. The map
// keys represent object attribute names.
func ObjectExact(value map[string]Check) objectExact {
	return objectExact{
		value: value,
	}
}

// createDeltaString prints the map keys that are present in mapA and not present in mapB
func createDeltaString[T any, V any](mapA map[string]T, mapB map[string]V, msgPrefix string) string {
	deltaMsg := ""

	deltaMap := make(map[string]T, len(mapA))
	maps.Copy(deltaMap, mapA)
	for key := range mapB {
		delete(deltaMap, key)
	}

	deltaKeys := slices.Sorted(maps.Keys(deltaMap))

	for i, k := range deltaKeys {
		if i == 0 {
			deltaMsg += msgPrefix
		} else if i != 0 {
			deltaMsg += ", "
		}
		deltaMsg += fmt.Sprintf("%q", k)
	}

	return deltaMsg
}
