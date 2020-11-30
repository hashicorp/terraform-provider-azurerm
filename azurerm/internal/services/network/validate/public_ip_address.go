package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
)

func PublicIPAddressID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.PublicIpAddressID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a resource id: %v", k, err))
		return
	}

	return warnings, errors
}
