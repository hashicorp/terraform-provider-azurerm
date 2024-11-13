// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
	"strconv"
)

var _ Check = tupleSizeExact{}

type tupleSizeExact struct {
	size int
}

// CheckValue verifies that the passed value is a tuple, map, object,
// or set, and contains a matching number of elements.
func (v tupleSizeExact) CheckValue(other any) error {
	otherVal, ok := other.([]any)

	if !ok {
		return fmt.Errorf("expected []any value for TupleSizeExact check, got: %T", other)
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

		return fmt.Errorf("expected %d %s for TupleSizeExact check, got %d %s", v.size, expectedElements, len(otherVal), actualElements)
	}

	return nil
}

// String returns the string representation of the value.
func (v tupleSizeExact) String() string {
	return strconv.FormatInt(int64(v.size), 10)
}

// TupleSizeExact returns a Check for asserting that
// a tuple has size elements.
func TupleSizeExact(size int) tupleSizeExact {
	return tupleSizeExact{
		size: size,
	}
}
