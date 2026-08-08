// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func DigitalTwinsTimeSeriesDatabaseConnectionName(v interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringLenBetween(3, 50),
		validation.StringDoesNotMatch(regexp.MustCompile(`^[0-9]+$`), "should not contain only numbers"),
		validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-]+[A-Za-z0-9]$`), "must begin with a letter or number, end with a letter or number and contain only letters, numbers, and hyphens"),
	)(v, k)
}
