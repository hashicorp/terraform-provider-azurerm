// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validation

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FloatBetween returns a SchemaValidateFunc which tests if the provided value
// is of type float64 and is between minVal and maxVal (inclusive).
func FloatBetween(minVal, maxVal float64) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(float64)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be float64", k))
			return
		}

		if v < minVal || v > maxVal {
			es = append(es, fmt.Errorf("expected %s to be in the range (%f - %f), got %f", k, minVal, maxVal, v))
			return
		}

		return
	}
}

// FloatAtLeast returns a SchemaValidateFunc which tests if the provided value
// is of type float and is at least minVal (inclusive)
func FloatAtLeast(minVal float64) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(float64)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be float", k))
			return
		}

		if v < minVal {
			es = append(es, fmt.Errorf("expected %s to be at least (%f), got %f", k, minVal, v))
			return
		}

		return
	}
}

// FloatAtMost returns a SchemaValidateFunc which tests if the provided value
// is of type float and is at most maxVal (inclusive)
func FloatAtMost(maxVal float64) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(float64)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be float", k))
			return
		}

		if v > maxVal {
			es = append(es, fmt.Errorf("expected %s to be at most (%f), got %f", k, maxVal, v))
			return
		}

		return
	}
}
