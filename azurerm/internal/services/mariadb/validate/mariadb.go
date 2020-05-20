package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mariadb/parse"
)

func MariaDbServerServerID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	if _, err := parse.MariaDbServerServerID(v); err != nil {
		errors = append(errors, fmt.Errorf("cannot parse %q as a MariaDB Server resource id: %v", k, err))
	}

	return warnings, errors
}

func MariaDbServerServerName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `^[0-9a-z][-0-9a-z]{1,61}[0-9a-z]$`); !m {
		return nil, append(regexErrs, fmt.Errorf("%q can contain only lowercase letters, numbers, and '-', but can't start or end with '-', and must be at least 3 characters and no more than 63 characters long.", k))
	}

	return nil, nil
}
