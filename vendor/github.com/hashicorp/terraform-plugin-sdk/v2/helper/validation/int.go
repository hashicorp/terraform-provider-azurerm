// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validation

import (
	"fmt"
	"math"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// IntBetween returns a SchemaValidateFunc which tests if the provided value
// is of type int and is between minVal and maxVal (inclusive)
func IntBetween(minVal, maxVal int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		v, ok := i.(int)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %s to be integer", k))
			return warnings, errors
		}

		if v < minVal || v > maxVal {
			errors = append(errors, fmt.Errorf("expected %s to be in the range (%d - %d), got %d", k, minVal, maxVal, v))
			return warnings, errors
		}

		return warnings, errors
	}
}

// IntAtLeast returns a SchemaValidateFunc which tests if the provided value
// is of type int and is at least minVal (inclusive)
func IntAtLeast(minVal int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		v, ok := i.(int)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %s to be integer", k))
			return warnings, errors
		}

		if v < minVal {
			errors = append(errors, fmt.Errorf("expected %s to be at least (%d), got %d", k, minVal, v))
			return warnings, errors
		}

		return warnings, errors
	}
}

// IntAtMost returns a SchemaValidateFunc which tests if the provided value
// is of type int and is at most maxVal (inclusive)
func IntAtMost(maxVal int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		v, ok := i.(int)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %s to be integer", k))
			return warnings, errors
		}

		if v > maxVal {
			errors = append(errors, fmt.Errorf("expected %s to be at most (%d), got %d", k, maxVal, v))
			return warnings, errors
		}

		return warnings, errors
	}
}

// IntDivisibleBy returns a SchemaValidateFunc which tests if the provided value
// is of type int and is divisible by a given number
func IntDivisibleBy(divisor int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		v, ok := i.(int)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %s to be integer", k))
			return warnings, errors
		}

		if math.Mod(float64(v), float64(divisor)) != 0 {
			errors = append(errors, fmt.Errorf("expected %s to be divisible by %d, got: %v", k, divisor, i))
			return warnings, errors
		}

		return warnings, errors
	}
}

// IntInSlice returns a SchemaValidateFunc which tests if the provided value
// is of type int and matches the value of an element in the valid slice
func IntInSlice(valid []int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		v, ok := i.(int)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %s to be integer", k))
			return warnings, errors
		}

		for _, validInt := range valid {
			if v == validInt {
				return warnings, errors
			}
		}

		errors = append(errors, fmt.Errorf("expected %s to be one of %v, got %d", k, valid, v))
		return warnings, errors
	}
}

// IntNotInSlice returns a SchemaValidateFunc which tests if the provided value
// is of type int and matches the value of an element in the valid slice
func IntNotInSlice(valid []int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		v, ok := i.(int)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %s to be integer", k))
			return warnings, errors
		}

		for _, validInt := range valid {
			if v == validInt {
				errors = append(errors, fmt.Errorf("expected %s to not be one of %v, got %d", k, valid, v))
			}
		}

		return warnings, errors
	}
}
