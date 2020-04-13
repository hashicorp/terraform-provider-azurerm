package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/parse"
)

func PolicyScopeID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.PolicyScopeID(v); err != nil {
		errors = append(errors, fmt.Errorf("can not parse %q as a Policy Scope ID: %+v", k, err))
		return
	}

	return warnings, errors
}
