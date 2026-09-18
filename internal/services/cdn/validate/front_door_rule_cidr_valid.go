// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// Evaluates if the passed CIDR is a valid IPv4 or IPv6 CIDR or not.
func FrontDoorRuleCidrIsValid(v interface{}, k string) ([]string, []error) {
	return validation.IsCIDR(v, k)
}
