package validate

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policyinsights/parse"
)

func RemediationName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	// The service returns error when name of remediation is too long
	// error: The remediation name cannot be empty and must not exceed '260' characters.
	// By my additional test, the name of remediation cannot contain the following characters: %^#/\&?.
	if len(v) == 0 || len(v) > 260 {
		errors = append(errors, fmt.Errorf("%s cannot be empty and must not exceed '260' characters", k))
		return
	}
	invalidCharacters := `%^#/\&?`
	if strings.ContainsAny(v, invalidCharacters) {
		errors = append(errors, fmt.Errorf("%s cannot contain the following characters: %s", k, invalidCharacters))
	}

	return warnings, errors
}

func RemediationID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.RemediationID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a resource id: %v", k, err))
		return
	}

	return warnings, errors
}

func RemediationScopeID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.RemediationScopeID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a resource id: %v", k, err))
		return
	}

	return warnings, errors
}
