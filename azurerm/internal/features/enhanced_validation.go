package features

import (
	"os"
	"strings"
)

// EnhancedValidationEnabled returns whether or not the feature for Enhanced Validation is
// enabled.
//
// This functionality calls out to the Azure MetaData Service to cache the list of supported
// Azure Locations for the specified Endpoint - and then uses that to provide enhanced validation
//
// This is enabled by default as of version 2.20 of the Azure Provider, and can be disabled by
// setting the Environment Variable `ARM_PROVIDER_ENHANCED_VALIDATION` to `false`.
func EnhancedValidationEnabled() bool {
	value := os.Getenv("ARM_PROVIDER_ENHANCED_VALIDATION")
	if value == "" {
		return true
	}

	return strings.EqualFold(value, "true")
}
