package validate

import (
	"fmt"
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/blueprint/parse"
)

func BlueprintDefinitionName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	// Portal says: Blueprint names can include letters, numbers or dashes. Spaces and other special characters are not allowed.
	// and the name can only have maximum 48 characters
	if matched := regexp.MustCompile(`^[A-Za-z0-9-]{1,48}$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%s can include letters, numbers or dashes. Spaces and other special characters are not allowed.", k))
	}
	return
}

func BlueprintDefinitionID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.BlueprintDefinitionID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a resource id: %v", k, err))
		return
	}

	return warnings, errors
}

func BlueprintDefinitionScopeID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.BlueprintDefinitionScopeID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a scope id: %v", k, err))
		return
	}

	return warnings, errors
}
