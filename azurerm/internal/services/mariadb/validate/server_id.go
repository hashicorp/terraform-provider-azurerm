package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mariadb/parse"
)

func ServerID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	if _, err := parse.ServerID(v); err != nil {
		errors = append(errors, fmt.Errorf("cannot parse %q as a MariaDB Server resource id: %v", k, err))
	}

	return warnings, errors
}
