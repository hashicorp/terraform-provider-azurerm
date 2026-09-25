// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// portal says: The name can only be an ASCII letter, digit, -, _, (, ), . and have a maximum length constraint of 90
func ManagementGroupName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9_().-]{1,90}$`), "can only consist of ASCII letters, digits, -, _, (, ), . , and cannot exceed the maximum length of 90")(i, k)
}
