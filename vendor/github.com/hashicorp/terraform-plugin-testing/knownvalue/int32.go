// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"encoding/json"
	"fmt"
	"strconv"
)

var _ Check = int32Exact{}

type int32Exact struct {
	value int32
}

// CheckValue determines whether the passed value is of type int32, and
// contains a matching int32 value.
func (v int32Exact) CheckValue(other any) error {
	jsonNum, ok := other.(json.Number)

	if !ok {
		return fmt.Errorf("expected json.Number value for Int32Exact check, got: %T", other)
	}

	otherVal, err := strconv.ParseInt(string(jsonNum), 10, 32)

	if err != nil {
		return fmt.Errorf("expected json.Number to be parseable as int32 value for Int32Exact check: %s", err)
	}

	if int32(otherVal) != v.value {
		return fmt.Errorf("expected value %d for Int32Exact check, got: %d", v.value, otherVal)
	}

	return nil
}

// String returns the string representation of the int32 value.
func (v int32Exact) String() string {
	return strconv.FormatInt(int64(v.value), 10)
}

// Int32Exact returns a Check for asserting equality between the
// supplied int32 and the value passed to the CheckValue method.
func Int32Exact(value int32) int32Exact {
	return int32Exact{
		value: value,
	}
}
