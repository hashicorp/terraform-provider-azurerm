// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ApplicationVersion(v interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringMatch(regexp.MustCompile(`^[-._\da-zA-Z]+$`), "can contain any combination of alphanumeric characters, hyphens, underscores, and periods"),
		validation.StringLenBetween(1, 64),
	)(v, k)
}
