// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func MountPath(i interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringLenBetween(2, 255),
		validation.StringMatch(regexp.MustCompile(`^(?:\/(?:[a-zA-Z][a-zA-Z0-9]*))+$`), "is not valid, must match the regular expression ^(?:\\/(?:[a-zA-Z][a-zA-Z0-9]*))+$"),
	)(i, k)
}
