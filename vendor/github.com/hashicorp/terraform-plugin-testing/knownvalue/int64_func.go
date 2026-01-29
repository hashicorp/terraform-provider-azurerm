// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"encoding/json"
	"fmt"
	"strconv"
)

var _ Check = int64Func{}

type int64Func struct {
	checkFunc func(v int64) error
}

// CheckValue determines whether the passed value is of type int64, and
// returns no error from the provided check function
func (v int64Func) CheckValue(other any) error {
	jsonNum, ok := other.(json.Number)

	if !ok {
		return fmt.Errorf("expected json.Number value for Int64Func check, got: %T", other)
	}

	otherVal, err := strconv.ParseInt(string(jsonNum), 10, 64)
	if err != nil {
		return fmt.Errorf("expected json.Number to be parseable as int64 value for Int64Func check: %s", err)
	}

	return v.checkFunc(otherVal)
}

// String returns the int64 representation of the value.
func (v int64Func) String() string {
	// Validation is up the the implementer of the function, so there are no
	// int64 literal or regex comparers to print here
	return "Int64Func"
}

// Int64Func returns a Check for passing the int64 value in state
// to the provided check function
func Int64Func(fn func(v int64) error) int64Func {
	return int64Func{
		checkFunc: fn,
	}
}
