// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func VaultName(v interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9-]{3,24}$`), "may only contain alphanumeric characters and dashes and must be between 3-24 chars"),
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z].*[a-zA-Z0-9]$`), "must start with a letter and end with a letter or number"),
		validation.StringDoesNotMatch(regexp.MustCompile(`--`), "cannot contain consecutive hyphens (\"--\")"),
	)(v, k)
}
