// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"encoding/json"
	"fmt"
	"strconv"
)

var _ Check = float32Exact{}

type float32Exact struct {
	value float32
}

// CheckValue determines whether the passed value is of type float32, and
// contains a matching float32 value.
func (v float32Exact) CheckValue(other any) error {
	jsonNum, ok := other.(json.Number)

	if !ok {
		return fmt.Errorf("expected json.Number value for Float32Exact check, got: %T", other)
	}

	otherVal, err := strconv.ParseFloat(string(jsonNum), 32)

	if err != nil {
		return fmt.Errorf("expected json.Number to be parseable as float32 value for Float32Exact check: %s", err)
	}

	if float32(otherVal) != v.value {
		return fmt.Errorf("expected value %s for Float32Exact check, got: %s", v.String(), strconv.FormatFloat(otherVal, 'f', -1, 32))
	}

	return nil
}

// String returns the string representation of the float32 value.
func (v float32Exact) String() string {
	return strconv.FormatFloat(float64(v.value), 'f', -1, 32)
}

// Float32Exact returns a Check for asserting equality between the
// supplied float32 and the value passed to the CheckValue method.
func Float32Exact(value float32) float32Exact {
	return float32Exact{
		value: value,
	}
}
