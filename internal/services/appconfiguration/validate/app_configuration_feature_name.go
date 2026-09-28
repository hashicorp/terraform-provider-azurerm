// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func AppConfigurationFeatureName(input interface{}, key string) ([]string, []error) {
	return validation.All(
		validation.StringDoesNotContainAny("%:"),
		validation.StringIsNotWhiteSpace,
	)(input, key)
}

func AppConfigurationFeatureKey(input interface{}, key string) ([]string, []error) {
	return validation.All(
		validation.StringDoesNotContainAny("%"),
		validation.StringIsNotWhiteSpace,
	)(input, key)
}
