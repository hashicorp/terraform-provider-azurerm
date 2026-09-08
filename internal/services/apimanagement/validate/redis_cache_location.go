// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func RedisCacheLocation(input interface{}, key string) (warnings []string, errors []error) {
	// "default" is valid in addition to any Azure location
	return validation.Any(
		validation.StringInSlice([]string{"default"}, false),
		location.EnhancedValidate,
	)(input, key)
}
