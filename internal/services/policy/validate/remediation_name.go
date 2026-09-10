// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func RemediationName(i interface{}, k string) (warnings []string, errors []error) {
	// The service returns error when name of remediation is too long
	// error: The remediation name cannot be empty and must not exceed '260' characters.
	// By additional testing, the name of remediation cannot contain the following characters: %^#/\&?.
	// Despite the service accepting remediation names with capitalized characters, in the response
	// all upper case characters will be converted to lower case. Therefore upper case letters are forbidden here.
	return validation.All(
		validation.StringLenBetween(1, 260),
		validation.StringDoesNotContainAny(`%^#/\&?`),
		validation.StringDoesNotMatch(regexp.MustCompile(`[A-Z]`), "cannot contain upper case letters"),
	)(i, k)
}
