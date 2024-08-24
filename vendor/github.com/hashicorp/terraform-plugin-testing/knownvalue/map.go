// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
	"sort"
)

var _ Check = mapExact{}

type mapExact struct {
	value map[string]Check
}

// CheckValue determines whether the passed value is of type map[string]any, and
// contains matching map entries.
func (v mapExact) CheckValue(other any) error {
	otherVal, ok := other.(map[string]any)

	if !ok {
		return fmt.Errorf("expected map[string]any value for MapExact check, got: %T", other)
	}

	if len(otherVal) != len(v.value) {
		expectedElements := "elements"
		actualElements := "elements"

		if len(v.value) == 1 {
			expectedElements = "element"
		}

		if len(otherVal) == 1 {
			actualElements = "element"
		}

		return fmt.Errorf("expected %d %s for MapExact check, got %d %s", len(v.value), expectedElements, len(otherVal), actualElements)
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
			return fmt.Errorf("missing element %s for MapExact check", k)
		}

		if err := v.value[k].CheckValue(otherValItem); err != nil {
			return fmt.Errorf("%s map element: %s", k, err)
		}
	}

	return nil
}

// String returns the string representation of the value.
func (v mapExact) String() string {
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

// MapExact returns a Check for asserting equality between the
// supplied map[string]Check and the value passed to the CheckValue method.
func MapExact(value map[string]Check) mapExact {
	return mapExact{
		value: value,
	}
}
