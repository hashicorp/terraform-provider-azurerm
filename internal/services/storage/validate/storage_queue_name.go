// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func StorageQueueName(v interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringMatch(regexp.MustCompile(`^[a-z0-9-]+$`), "only lowercase alphanumeric characters and hyphens allowed"),
		validation.StringDoesNotMatch(regexp.MustCompile(`^-`), "cannot start with a hyphen"),
		validation.StringDoesNotMatch(regexp.MustCompile(`-$`), "cannot end with a hyphen"),
		validation.StringLenBetween(3, 63),
	)(v, k)
}
