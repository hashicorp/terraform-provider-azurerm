// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var durableTaskNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]$`)

func DurableTaskName(i interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringLenBetween(3, 63),
		validation.StringMatch(durableTaskNameRegex, "must start and end with alphanumeric characters and can contain hyphens"),
	)(i, k)
}
