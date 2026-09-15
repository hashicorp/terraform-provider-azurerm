// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func StringDoesNotContainUpperCaseLetter(input interface{}, k string) ([]string, []error) {
	return validation.StringDoesNotMatch(regexp.MustCompile(`[\p{Lu}\p{Lt}]`), "expected value to not contain any uppercase letter")(input, k)
}
