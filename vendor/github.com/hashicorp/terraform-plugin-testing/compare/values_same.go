// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compare

import (
	"fmt"
	"reflect"
)

var _ ValueComparer = valuesSame{}

type valuesSame struct{}

// CompareValues determines whether each value in the sequence of the supplied values
// is the same as the preceding value.
func (v valuesSame) CompareValues(values ...any) error {
	for i := 1; i < len(values); i++ {
		if !reflect.DeepEqual(values[i-1], values[i]) {
			return fmt.Errorf("expected values to be the same, but they differ: %v != %v", values[i-1], values[i])
		}
	}

	return nil
}

// ValuesSame returns a ValueComparer for asserting that each value in the sequence of
// the values supplied to the CompareValues method is the same as the preceding value.
func ValuesSame() valuesSame {
	return valuesSame{}
}
