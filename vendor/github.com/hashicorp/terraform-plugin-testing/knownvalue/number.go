// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package knownvalue

import (
	"encoding/json"
	"fmt"
	"math/big"
)

var _ Check = numberExact{}

type numberExact struct {
	value *big.Float
}

// CheckValue determines whether the passed value is of type *big.Float, and
// contains a matching *big.Float value.
func (v numberExact) CheckValue(other any) error {
	if v.value == nil {
		return fmt.Errorf("value in NumberExact check is nil")
	}

	jsonNum, ok := other.(json.Number)

	if !ok {
		return fmt.Errorf("expected json.Number value for NumberExact check, got: %T", other)
	}

	otherVal, _, err := big.ParseFloat(jsonNum.String(), 10, 512, big.ToNearestEven)

	if err != nil {
		return fmt.Errorf("expected json.Number to be parseable as big.Float value for NumberExact check: %s", err)
	}

	if v.value.Cmp(otherVal) != 0 {
		return fmt.Errorf("expected value %s for NumberExact check, got: %s", v.String(), otherVal.Text('f', -1))
	}

	return nil
}

// String returns the string representation of the *big.Float value.
func (v numberExact) String() string {
	return v.value.Text('f', -1)
}

// NumberExact returns a Check for asserting equality between the
// supplied *big.Float and the value passed to the CheckValue method.
// The CheckValue method uses 512-bit precision to perform this assertion.
func NumberExact(value *big.Float) numberExact {
	return numberExact{
		value: value,
	}
}
