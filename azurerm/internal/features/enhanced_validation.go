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
// This can be enabled using the Environment Variable `ARM_PROVIDER_ENHANCED_VALIDATION` and
// defaults to 'false' at the present time - but may change in a future release.
func EnhancedValidationEnabled() bool {
	return strings.EqualFold(os.Getenv("ARM_PROVIDER_ENHANCED_VALIDATION"), "true")
}
