package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/logic/parse"
)

func IntegrationServiceEnvironmentID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	if _, err := parse.IntegrationServiceEnvironmentID(v); err != nil {
		errors = append(errors, fmt.Errorf("cannot parse %q as an Integration Service Environment ID: %+v", k, err))
	}

	return warnings, errors
}
