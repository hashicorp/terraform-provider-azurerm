package validate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/synapse/parse"
)

func SynapseWorkspaceName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	// The name attribute rules are :
	// 1. can contain only lowercase letters or numbers.
	// 2. must not end with the string 'ondemand'
	// 3. The value must be between 1 and 45 characters long

	if !regexp.MustCompile(`^[a-z\d]{1,45}$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%s can contain only lowercase letters or numbers, and be between 1 and 45 characters long", k))
		return
	}
	if strings.HasSuffix(v, "ondemand") {
		errors = append(errors, fmt.Errorf("%s must not end with the string 'ondemand'", k))
		return
	}

	return warnings, errors
}

func SqlAdministratorLoginName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	// The name attribute rules are :
	// 1. can contain only letters or numbers.
	// 2. must start with letter
	// 3. The value must be between 1 and 128 characters long

	if !regexp.MustCompile(`^[a-zA-Z][a-zA-Z\d]{0,127}$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%s can contain only letters or numbers, must start with a letter, The value must be between 1 and 128 characters long", k))
		return
	}

	return warnings, errors
}

func SynapseWorkspaceID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.SynapseWorkspaceID(v); err != nil {
		errors = append(errors, fmt.Errorf("can not parse %q as a synapse workspace resource id: %v", k, err))
		return
	}

	return warnings, errors
}
