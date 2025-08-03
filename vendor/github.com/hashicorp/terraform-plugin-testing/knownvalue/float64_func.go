// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"encoding/json"
	"fmt"
	"strconv"
)

var _ Check = float64Func{}

type float64Func struct {
	checkFunc func(v float64) error
}

// CheckValue determines whether the passed value is of type float64, and
// returns no error from the provided check function
func (v float64Func) CheckValue(other any) error {
	jsonNum, ok := other.(json.Number)

	if !ok {
		return fmt.Errorf("expected json.Number value for Float64Func check, got: %T", other)
	}

	otherVal, err := strconv.ParseFloat(string(jsonNum), 64)
	if err != nil {
		return fmt.Errorf("expected json.Number to be parseable as float64 value for Float64Func check: %s", err)
	}

	return v.checkFunc(otherVal)
}

// String returns the float64 representation of the value.
func (v float64Func) String() string {
	// Validation is up the the implementer of the function, so there are no
	// float64 literal or regex comparers to print here
	return "Float64Func"
}

// Float64Func returns a Check for passing the float64 value in state
// to the provided check function
func Float64Func(fn func(v float64) error) float64Func {
	return float64Func{
		checkFunc: fn,
	}
}
