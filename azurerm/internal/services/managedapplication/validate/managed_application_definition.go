package validate

import (
	"fmt"
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managedapplication/parse"
)

func ManagedApplicationDefinitionID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.ManagedApplicationDefinitionID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a resource id: %v", k, err))
		return
	}

	return warnings, errors
}

func ManagedApplicationDefinitionName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\da-zA-Z]{3,64}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 3 and 64 characters in length and contains only letters or numbers.", k))
	}

	return warnings, errors
}

func ManagedApplicationDefinitionDisplayName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\s\S]{4,60}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 4 and 60 characters in length.", k))
	}

	return warnings, errors
}

func ManagedApplicationDefinitionDescription(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\s\S]{0,200}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q should not exceed 200 characters in length.", k))
	}

	return warnings, errors
}
