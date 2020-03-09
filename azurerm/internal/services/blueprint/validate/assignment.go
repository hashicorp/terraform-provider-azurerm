package validate

import (
	"fmt"
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/blueprint/parse"
)

func BlueprintAssignmentName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	// Portal says: Blueprint names can include letters, numbers or dashes. Spaces and other special characters are not allowed.
	if matched := regexp.MustCompile(`^[A-Za-z0-9-]{1,90}$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%s can include letters, numbers or dashes, spaces and other special characters are not allowed. %s cannot exceed the maximum length of 90", k, k))
	}
	return
}

func BlueprintAssignmentID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.BlueprintAssignmentID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a resource id: %v", k, err))
		return
	}

	return warnings, errors
}

func BlueprintAssignmentScopeID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.BlueprintAssignmentScopeID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a scope id: %v", k, err))
		return
	}

	return warnings, errors
}
