package validate

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/logic/parse"
)

func IntegrationAccountName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^[\w-().]{1,80}$`), `Integration name can contain only letters, numbers, '_','-', '(', ')' or '.'`,
	)
}

func IntegrationAccountID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	if _, err := parse.IntegrationAccountID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a Integration Account resource id: %v", k, err))
	}

	return warnings, errors
}

// https://docs.microsoft.com/en-us/azure/azure-resource-manager/management/resource-name-rules
// 1-80	Alphanumerics, hyphens, periods, and underscores.
func IntegrationServiceEnvironmentName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^[\w-.]{1,80}$`), `Integration Service Environment names must be between 1 and 80 characters in length, contain only letters, numbers, underscores, hyphens and periods.`,
	)
}

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
