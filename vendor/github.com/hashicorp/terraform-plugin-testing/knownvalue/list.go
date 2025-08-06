// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"fmt"
)

var _ Check = listExact{}

type listExact struct {
	value []Check
}

// CheckValue determines whether the passed value is of type []any, and
// contains matching slice entries in the same sequence.
func (v listExact) CheckValue(other any) error {
	otherVal, ok := other.([]any)

	if !ok {
		return fmt.Errorf("expected []any value for ListExact check, got: %T", other)
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

		return fmt.Errorf("expected %d %s for ListExact check, got %d %s", len(v.value), expectedElements, len(otherVal), actualElements)
	}

	for i := 0; i < len(v.value); i++ {
		if err := v.value[i].CheckValue(otherVal[i]); err != nil {
			return fmt.Errorf("list element index %d: %s", i, err)
		}
	}

	return nil
}

// String returns the string representation of the value.
func (v listExact) String() string {
	var listVals []string

	for _, val := range v.value {
		listVals = append(listVals, val.String())
	}

	return fmt.Sprintf("%s", listVals)
}

// ListExact returns a Check for asserting equality between the
// supplied []Check and the value passed to the CheckValue method.
// This is an order-dependent check.
func ListExact(value []Check) listExact {
	return listExact{
		value: value,
	}
}
