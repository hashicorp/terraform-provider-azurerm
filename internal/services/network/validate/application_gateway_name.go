// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ApplicationGatewayName(i interface{}, k string) ([]string, []error) {
	// Validate name: 1-80 chars, begin with letter/number, end with letter/number/underscore,
	// and may contain letters, numbers, underscores, periods, or hyphens in the middle
	return validation.StringMatch(
		regexp.MustCompile(`^[a-zA-Z\d]$|^[a-zA-Z\d][a-zA-Z\d-_.]{0,78}[a-zA-Z\d_]$`),
		"must be between 1 - 80 characters long, begin with a letter or number, end with a letter, number or underscore (_), and may contain only alphanumeric characters, underscores (_), hyphens (-), and periods (.)",
	)(i, k)
}
