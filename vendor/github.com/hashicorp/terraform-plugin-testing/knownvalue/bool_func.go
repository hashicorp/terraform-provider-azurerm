// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import "fmt"

var _ Check = boolFunc{}

type boolFunc struct {
	checkFunc func(v bool) error
}

// CheckValue determines whether the passed value is of type bool, and
// returns no error from the provided check function
func (v boolFunc) CheckValue(other any) error {
	val, ok := other.(bool)

	if !ok {
		return fmt.Errorf("expected bool value for BoolFunc check, got: %T", other)
	}

	return v.checkFunc(val)
}

// String returns the bool representation of the value.
func (v boolFunc) String() string {
	// Validation is up the the implementer of the function, so there are no
	// bool literal or regex comparers to print here
	return "BoolFunc"
}

// BoolFunc returns a Check for passing the bool value in state
// to the provided check function
func BoolFunc(fn func(v bool) error) boolFunc {
	return boolFunc{
		checkFunc: fn,
	}
}
