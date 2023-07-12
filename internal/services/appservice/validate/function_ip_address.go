// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// IsIpOrCIDRRange is a SchemaValidateFunc which tests if the provided value is of type string and a valid IP/CIDR range
func IsIpOrCIDRRangeList(i interface{}, k string) ([]string, []error) {
	var allWarnings []string
	var allErrors []error
	v, ok := i.(string)
	if !ok {
		allErrors = append(allErrors, fmt.Errorf("expected type of %s to be string", k))
		return allWarnings, allErrors
	}
	validator := validation.Any(validation.IsCIDR, validation.IsIPAddress)
	for _, elem := range strings.Split(v, ",") {
		warnings, errors := validator(elem, k)
		if len(warnings) > 0 {
			allWarnings = append(allWarnings, warnings...)
		}
		if len(errors) > 0 {
			allErrors = append(allErrors, errors...)
		}
	}
	return allWarnings, allErrors
}
