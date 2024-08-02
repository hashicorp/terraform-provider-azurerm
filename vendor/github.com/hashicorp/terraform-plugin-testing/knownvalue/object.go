// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
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
		expectedAttributes := "attributes"
		actualAttributes := "attributes"

		if len(v.value) == 1 {
			expectedAttributes = "attribute"
		}

		if len(otherVal) == 1 {
			actualAttributes = "attribute"
		}

		return fmt.Errorf("expected %d %s for ObjectExact check, got %d %s", len(v.value), expectedAttributes, len(otherVal), actualAttributes)
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
