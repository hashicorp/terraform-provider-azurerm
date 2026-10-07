// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package features

import (
	"os"
	"strings"
)

// EnhancedValidationLocationsEnabled returns whether Enhanced Validation for Locations is enabled.
//
// This checks the `ARM_PROVIDER_ENHANCED_VALIDATION_LOCATIONS` environment variable first, before falling back to the default.
func EnhancedValidationLocationsEnabled() bool {
	value := os.Getenv("ARM_PROVIDER_ENHANCED_VALIDATION_LOCATIONS")
	if value != "" {
		return strings.EqualFold(value, "true")
	}

	return false
}

// EnhancedValidationResourceProvidersEnabled returns whether Enhanced Validation for Resource Providers is enabled.
//
// This checks the `ARM_PROVIDER_ENHANCED_VALIDATION_RESOURCE_PROVIDERS` environment variable first, before falling back to the default.
func EnhancedValidationResourceProvidersEnabled() bool {
	value := os.Getenv("ARM_PROVIDER_ENHANCED_VALIDATION_RESOURCE_PROVIDERS")
	if value != "" {
		return strings.EqualFold(value, "true")
	}

	return false
}

// EnhancedValidationPreflightEnabled returns whether Azure Preflight Validation is enabled.
//
// Preflight validation is always opt-in and defaults to false in all provider versions.
// Set ARM_PROVIDER_ENHANCED_VALIDATION_PREFLIGHT_ENABLED=true to enable it without
// requiring an explicit config block.
func EnhancedValidationPreflightEnabled() bool {
	return strings.EqualFold(os.Getenv("ARM_PROVIDER_ENHANCED_VALIDATION_PREFLIGHT_ENABLED"), "true")
}

// EnhancedValidationLocationFallback returns the location fallback string for Azure Preflight Validation.
//
// This checks the `ARM_PROVIDER_ENHANCED_VALIDATION_LOCATION_FALLBACK` environment variable,
// which can be used to set a default location for resources that don't have a specific location.
func EnhancedValidationLocationFallback() *string {
	value := os.Getenv("ARM_PROVIDER_ENHANCED_VALIDATION_LOCATION_FALLBACK")
	if value != "" {
		return &value
	}
	return nil
}
