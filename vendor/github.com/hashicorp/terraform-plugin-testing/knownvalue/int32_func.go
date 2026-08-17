// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"encoding/json"
	"fmt"
	"strconv"
)

var _ Check = int32Func{}

type int32Func struct {
	checkFunc func(v int32) error
}

// CheckValue determines whether the passed value is of type int32, and
// returns no error from the provided check function
func (v int32Func) CheckValue(other any) error {
	jsonNum, ok := other.(json.Number)

	if !ok {
		return fmt.Errorf("expected json.Number value for Int32Func check, got: %T", other)
	}

	otherVal, err := strconv.ParseInt(string(jsonNum), 10, 32)
	if err != nil {
		return fmt.Errorf("expected json.Number to be parseable as int32 value for Int32Func check: %s", err)
	}

	return v.checkFunc(int32(otherVal))
}

// String returns the int32 representation of the value.
func (v int32Func) String() string {
	// Validation is up the the implementer of the function, so there are no
	// int32 literal or regex comparers to print here
	return "Int32Func"
}

// Int32Func returns a Check for passing the int32 value in state
// to the provided check function
func Int32Func(fn func(v int32) error) int32Func {
	return int32Func{
		checkFunc: fn,
	}
}
