// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"encoding/json"
	"fmt"
	"strconv"
)

var _ Check = float32Func{}

type float32Func struct {
	checkFunc func(v float32) error
}

// CheckValue determines whether the passed value is of type float32, and
// returns no error from the provided check function
func (v float32Func) CheckValue(other any) error {
	jsonNum, ok := other.(json.Number)

	if !ok {
		return fmt.Errorf("expected json.Number value for Float32Func check, got: %T", other)
	}

	otherVal, err := strconv.ParseFloat(string(jsonNum), 32)
	if err != nil {
		return fmt.Errorf("expected json.Number to be parseable as float32 value for Float32Func check: %s", err)
	}

	return v.checkFunc(float32(otherVal))
}

// String returns the float32 representation of the value.
func (v float32Func) String() string {
	// Validation is up the the implementer of the function, so there are no
	// float32 literal or regex comparers to print here
	return "Float32Func"
}

// Float32Func returns a Check for passing the float32 value in state
// to the provided check function
func Float32Func(fn func(v float32) error) float32Func {
	return float32Func{
		checkFunc: fn,
	}
}
