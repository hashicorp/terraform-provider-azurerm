package validate

import (
	"fmt"
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/parse"
)

// TODO -- remove this function when validation function is implemented in azurerm_policy_definition
func PolicyDefinitionID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	regex := regexp.MustCompile(`/providers/[Mm]icrosoft\.[Aa]uthorization/policy[Dd]efinitions/`)
	segments := regex.Split(v, -1)

	if len(segments) != 2 {
		errors = append(errors, fmt.Errorf("expected %q to have 2 segments after splition", k))
		return
	}

	scope := segments[0]
	// scope should be a resource group ID, or Management Group ID
	if _, err := parse.PolicyScopeID(scope); err != nil {
		errors = append(errors, fmt.Errorf("cannot parse %q as valid scope id: %+v", k, err))
		return
	}

	name := segments[1]
	// policy definition should have a name
	if name == "" {
		errors = append(errors, fmt.Errorf("expected %q to have a non-empty name", k))
	}

	return warnings, errors
}
