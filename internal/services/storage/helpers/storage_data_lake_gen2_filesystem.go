// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ValidateStorageDataLakeGen2FileSystemName(v any, k string) ([]string, []error) {
	return validation.All(
		validation.StringMatch(regexp.MustCompile(`^\$root$|^[0-9a-z-]+$`), "only lowercase alphanumeric characters and hyphens allowed"),
		validation.StringLenBetween(3, 63),
		validation.StringDoesNotMatch(regexp.MustCompile(`^-`), "cannot begin with a hyphen"),
	)(v, k)
}
