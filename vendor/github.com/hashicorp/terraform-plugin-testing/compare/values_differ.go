// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compare

import (
	"fmt"
	"reflect"
)

var _ ValueComparer = valuesDiffer{}

type valuesDiffer struct{}

// CompareValues determines whether each value in the sequence of the supplied values
// differs from the preceding value.
func (v valuesDiffer) CompareValues(values ...any) error {
	for i := 1; i < len(values); i++ {
		if reflect.DeepEqual(values[i-1], values[i]) {
			return fmt.Errorf("expected values to differ, but they are the same: %v == %v", values[i-1], values[i])
		}
	}

	return nil
}

// ValuesDiffer returns a ValueComparer for asserting that each value in the sequence of
// the values supplied to the CompareValues method differs from the preceding value.
func ValuesDiffer() valuesDiffer {
	return valuesDiffer{}
}
