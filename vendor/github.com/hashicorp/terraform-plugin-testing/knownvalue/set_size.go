// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
	"strconv"
)

var _ Check = setSizeExact{}

type setSizeExact struct {
	size int
}

// CheckValue verifies that the passed value is a list, map, object,
// or set, and contains a matching number of elements.
func (v setSizeExact) CheckValue(other any) error {
	otherVal, ok := other.([]any)

	if !ok {
		return fmt.Errorf("expected []any value for SetElementExact check, got: %T", other)
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

		return fmt.Errorf("expected %d %s for SetElementExact check, got %d %s", v.size, expectedElements, len(otherVal), actualElements)
	}

	return nil
}

// String returns the string representation of the value.
func (v setSizeExact) String() string {
	return strconv.FormatInt(int64(v.size), 10)
}

// SetSizeExact returns a Check for asserting that
// a set has size elements.
func SetSizeExact(size int) setSizeExact {
	return setSizeExact{
		size: size,
	}
}
