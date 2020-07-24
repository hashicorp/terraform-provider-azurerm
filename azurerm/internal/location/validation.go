package location

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

// this is only here to aid testing
var enhancedEnabled = features.EnhancedValidationEnabled()

// EnhancedValidate returns a validation function which attempts to validate the location
// against the list of Locations supported by this Azure Location.
//
// NOTE: this is best-effort - if the users offline, or the API doesn't return it we'll
// fall back to the original approach
func EnhancedValidate(i interface{}, k string) ([]string, []error) {
	if !enhancedEnabled || supportedLocations == nil {
		return validate.NoEmptyStrings(i, k)
	}

	return enhancedValidation(i, k)
}

func enhancedValidation(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	normalizedUserInput := Normalize(v)
	if normalizedUserInput == "" {
		return nil, []error{fmt.Errorf("%q must not be empty", k)}
	}

	// supportedLocations can be nil if the users offline
	if supportedLocations != nil {
		found := false
		for _, loc := range *supportedLocations {
			if normalizedUserInput == Normalize(loc) {
				found = true
				break
			}
		}

		if !found {
			locations := strings.Join(*supportedLocations, ",")
			return nil, []error{
				fmt.Errorf("%q was not found in the list of supported Azure Locations: %q", normalizedUserInput, locations),
			}
		}
	}

	return nil, nil
}
