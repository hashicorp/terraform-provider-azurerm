// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import "fmt"

var _ Check = stringExact{}

type stringExact struct {
	value string
}

// CheckValue determines whether the passed value is of type string, and
// contains a matching sequence of bytes.
func (v stringExact) CheckValue(other any) error {
	otherVal, ok := other.(string)

	if !ok {
		return fmt.Errorf("expected string value for StringExact check, got: %T", other)
	}

	if otherVal != v.value {
		return fmt.Errorf("expected value %s for StringExact check, got: %s", v.value, otherVal)
	}

	return nil
}

// String returns the string representation of the value.
func (v stringExact) String() string {
	return v.value
}

// StringExact returns a Check for asserting equality between the
// supplied string and a value passed to the CheckValue method.
func StringExact(value string) stringExact {
	return stringExact{
		value: value,
	}
}
