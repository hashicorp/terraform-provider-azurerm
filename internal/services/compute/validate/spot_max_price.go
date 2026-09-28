// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// SpotMaxPrice validates the price provided is a valid Spot Price for the Compute
// API (and downstream API's which use this like AKS) - either -1 (the current VM price)
// or at least 0.00001
func SpotMaxPrice(i interface{}, k string) ([]string, []error) {
	return validation.Any(
		validation.FloatInSlice([]float64{-1}),
		validation.FloatAtLeast(0.00001),
	)(i, k)
}
