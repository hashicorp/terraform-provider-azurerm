// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func DatabaseAutoPauseDelay(i interface{}, k string) (warnings []string, errors []error) {
	// -1 (disabled) is valid in addition to the range
	return validation.Any(
		validation.IntInSlice([]int{-1}),
		validation.IntBetween(15, 10080),
	)(i, k)
}
