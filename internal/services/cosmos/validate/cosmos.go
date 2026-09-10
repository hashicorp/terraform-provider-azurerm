// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// Portal: The value must contain only alphanumeric characters or the following: -
func CosmosAccountName(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile("^[-a-z0-9]{3,50}$"), "must be 3 - 50 characters long, and contain only letters, numbers and hyphens")(v, k)
}

func CosmosEntityName(v interface{}, k string) ([]string, []error) {
	return validation.StringLenBetween(1, 255)(v, k)
}

func CosmosThroughput(v interface{}, k string) ([]string, []error) {
	return validation.All(validation.IntAtLeast(400), validation.IntDivisibleBy(100))(v, k)
}

func CosmosMaxThroughput(i interface{}, k string) ([]string, []error) {
	return validation.All(validation.IntAtLeast(1000), validation.IntDivisibleBy(1000))(i, k)
}
