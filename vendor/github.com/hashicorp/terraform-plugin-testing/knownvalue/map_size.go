// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
	"strconv"
)

var _ Check = mapSizeExact{}

type mapSizeExact struct {
	size int
}

// CheckValue verifies that the passed value is a list, map, object,
// or set, and contains a matching number of elements.
func (v mapSizeExact) CheckValue(other any) error {
	otherVal, ok := other.(map[string]any)

	if !ok {
		return fmt.Errorf("expected map[string]any value for MapSizeExact check, got: %T", other)
	}

	if len(otherVal) != v.size {
		expectedElements := "elements"
		actualElements := "elements"

		if v.size == 1 {
			expectedElements = "element"
		}

		if len(otherVal) == 1 {
			actualElements = "element"
		}

		return fmt.Errorf("expected %d %s for MapSizeExact check, got %d %s", v.size, expectedElements, len(otherVal), actualElements)
	}

	return nil
}

// String returns the string representation of the value.
func (v mapSizeExact) String() string {
	return strconv.Itoa(v.size)
}

// MapSizeExact returns a Check for asserting that
// a map has size elements.
func MapSizeExact(size int) mapSizeExact {
	return mapSizeExact{
		size: size,
	}
}
