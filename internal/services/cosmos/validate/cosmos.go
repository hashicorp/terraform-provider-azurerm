// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func CosmosAccountName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	// Portal: The value must contain only alphanumeric characters or the following: -
	if matched := regexp.MustCompile("^[-a-z0-9]{3,50}$").Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%s name must be 3 - 50 characters long, contain only letters, numbers and hyphens", k))
	}

	return warnings, errors
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
