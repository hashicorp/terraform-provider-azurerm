package validate

import (
	"fmt"
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/synapse/parse"
)

func SqlPoolName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	// The name attribute rules are :
	// 1. can contain only letters, numbers and underscore.
	// 2. The value must be between 1 and 15 characters long

	if !regexp.MustCompile(`^[a-zA-Z_\d]{1,15}$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%s can contain only letters, numbers or underscore, The value must be between 1 and 15 characters long", k))
		return
	}

	return warnings, errors
}

func SqlPoolID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.SynapseSqlPoolID(v); err != nil {
		errors = append(errors, fmt.Errorf("can not parse %q as a Synapse Sql Pool resource id: %v", k, err))
		return
	}

	return warnings, errors
}
