// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import "fmt"

var _ Check = stringFunc{}

type stringFunc struct {
	checkFunc func(v string) error
}

// CheckValue determines whether the passed value is of type string, and
// returns no error from the provided check function
func (v stringFunc) CheckValue(value any) error {
	val, ok := value.(string)

	if !ok {
		return fmt.Errorf("expected string value for StringFunc check, got: %T", value)
	}

	return v.checkFunc(val)
}

// String returns the string representation of the value.
func (v stringFunc) String() string {
	// Validation is up the the implementer of the function, so there are no
	// string literal or regex comparers to print here
	return "StringFunc"
}

// StringFunc returns a Check for passing the string value in state
// to the provided check function
func StringFunc(fn func(v string) error) stringFunc {
	return stringFunc{
		checkFunc: fn,
	}
}
