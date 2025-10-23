// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"encoding/json"
	"fmt"
	"math/big"
)

var _ Check = numberFunc{}

type numberFunc struct {
	checkFunc func(v *big.Float) error
}

// CheckValue determines whether the passed value is of type int64, and
// returns no error from the provided check function
func (v numberFunc) CheckValue(other any) error {
	jsonNum, ok := other.(json.Number)

	if !ok {
		return fmt.Errorf("expected json.Number value for NumberFunc check, got: %T", other)
	}

	otherVal, _, err := big.ParseFloat(jsonNum.String(), 10, 512, big.ToNearestEven)
	if err != nil {
		return fmt.Errorf("expected json.Number to be parseable as big.Float value for NumberFunc check: %s", err)
	}

	return v.checkFunc(otherVal)
}

// String returns the int64 representation of the value.
func (v numberFunc) String() string {
	// Validation is up the the implementer of the function, so there are no
	// int64 literal or regex comparers to print here
	return "NumberFunc"
}

// NumberFunc returns a Check for passing the int64 value in state
// to the provided check function
func NumberFunc(fn func(v *big.Float) error) numberFunc {
	return numberFunc{
		checkFunc: fn,
	}
}
