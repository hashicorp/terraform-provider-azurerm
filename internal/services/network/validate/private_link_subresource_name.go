// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func PrivateLinkSubResourceName(i interface{}, k string) (_ []string, errors []error) {
	// empty is valid; otherwise the name must begin and end with an alphanumeric character, be between
	// 3 and 63 characters in length and only contain letters, numbers, underscores, periods, dashes, and spaces
	return validation.Any(
		validation.StringIsEmpty,
		validation.StringMatch(regexp.MustCompile(`^([a-zA-Z0-9])([\w .-]{1,61})([a-zA-Z0-9])$`), "must begin and end with an alphanumeric character, be between 3 and 63 characters in length, only contain letters, numbers, underscores, periods, dashes, and spaces"),
	)(i, k)
}
