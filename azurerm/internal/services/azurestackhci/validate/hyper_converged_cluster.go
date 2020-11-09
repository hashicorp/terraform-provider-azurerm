package validate

import (
	"fmt"
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/azurestackhci/parse"
)

func HyperConvergedClusterID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.HyperConvergedClusterID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a resource id: %v", k, err))
		return
	}

	return warnings, errors
}

func HyperConvergedClusterName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if matched := regexp.MustCompile(`^.{1,260}$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q may only be up to 260 characters in length", k))
	}

	return warnings, errors
}
