// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package features

import (
	"fmt"
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

// EnhancedValidationLocationsEnabled returns whether Enhanced Validation for Locations is enabled.
//
// This checks the `ARM_PROVIDER_ENHANCED_VALIDATION_LOCATIONS` environment variable first,
// falling back to the legacy `ARM_PROVIDER_ENHANCED_VALIDATION` environment variable, then to the
// version default (enabled in 4.x, disabled in 5.0).
func EnhancedValidationLocationsEnabled() bool {
	// Check the locations-specific env var first
	value := os.Getenv("ARM_PROVIDER_ENHANCED_VALIDATION_LOCATIONS")
	if value != "" {
		return strings.EqualFold(value, "true")
	}

	// In 4.x, fall back to the legacy environment variable
	if !FivePointOh() {
		return EnhancedValidationEnabled()
	}

	// In 5.0, default to disabled
	return false
}

// EnhancedValidationResourceProvidersEnabled returns whether Enhanced Validation for Resource Providers is enabled.
//
// This checks the `ARM_PROVIDER_ENHANCED_VALIDATION_RESOURCE_PROVIDERS` environment variable first,
// falling back to the legacy `ARM_PROVIDER_ENHANCED_VALIDATION` environment variable, then to the
// version default (enabled in 4.x, disabled in 5.0).
func EnhancedValidationResourceProvidersEnabled() bool {
	// Check the resource-providers-specific env var first
	value := os.Getenv("ARM_PROVIDER_ENHANCED_VALIDATION_RESOURCE_PROVIDERS")
	if value != "" {
		return strings.EqualFold(value, "true")
	}

	// In 4.x, fall back to the legacy environment variable
	if !FivePointOh() {
		return EnhancedValidationEnabled()
	}

	// In 5.0, default to disabled
	return false
}

// ValidateEnhancedValidationEnvVars validates the enhanced validation environment variables.
//
// In version 5.0, the legacy `ARM_PROVIDER_ENHANCED_VALIDATION` environment variable has been
// removed - an error is returned if it is set, directing users to migrate to the specific
// environment variables or the `enhanced_validation` provider block.
//
// In version 4.x, the legacy environment variable is still supported, but it cannot be set
// at the same time as any of the specific environment variables.
func ValidateEnhancedValidationEnvVars() error {
	legacyEnv := os.Getenv("ARM_PROVIDER_ENHANCED_VALIDATION")
	if legacyEnv == "" {
		return nil
	}

	// In 5.0, the legacy env var is no longer supported
	if FivePointOh() {
		return fmt.Errorf("the environment variable `ARM_PROVIDER_ENHANCED_VALIDATION` has been removed in v5.0 of the AzureRM Provider - please use the `enhanced_validation` provider block or the specific environment variables `ARM_PROVIDER_ENHANCED_VALIDATION_LOCATIONS` and `ARM_PROVIDER_ENHANCED_VALIDATION_RESOURCE_PROVIDERS` instead")
	}

	// In 4.x, check for conflicts with specific env vars
	var conflicts []string
	if v := os.Getenv("ARM_PROVIDER_ENHANCED_VALIDATION_LOCATIONS"); v != "" {
		conflicts = append(conflicts, "ARM_PROVIDER_ENHANCED_VALIDATION_LOCATIONS")
	}
	if v := os.Getenv("ARM_PROVIDER_ENHANCED_VALIDATION_RESOURCE_PROVIDERS"); v != "" {
		conflicts = append(conflicts, "ARM_PROVIDER_ENHANCED_VALIDATION_RESOURCE_PROVIDERS")
	}

	if len(conflicts) > 0 {
		return fmt.Errorf("the environment variable `ARM_PROVIDER_ENHANCED_VALIDATION` cannot be set at the same time as %v - please either use the legacy `ARM_PROVIDER_ENHANCED_VALIDATION` or the specific environment variables, but not both", conflicts)
	}

	return nil
}
