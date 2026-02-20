// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package features

import (
	"os"
	"strings"
)

// EnhancedValidationEnabled returns whether the feature for Enhanced Validation is enabled.
//
// This functionality calls out to the Azure MetaData Service to cache the list of supported
// Azure Locations and Resource Providers for the specified Endpoint - and then uses that to
// provide enhanced validation. When enabled, invalid locations or resource providers are caught
// at `terraform plan` time. When disabled, these errors are caught at `terraform apply` time
// when Azure rejects the request.
//
// This is enabled by default in version 4.x and disabled by default as of version 5.0 of the
// Azure Provider. The default can be overridden by setting the Environment Variable
// `ARM_PROVIDER_ENHANCED_VALIDATION` to `true` or `false`.
func EnhancedValidationEnabled() bool {
	value := os.Getenv("ARM_PROVIDER_ENHANCED_VALIDATION")
	if value == "" {
		// In 5.0, default to disabled; in 4.x, default to enabled
		return !FivePointOh()
	}

	return strings.EqualFold(value, "true")
}
