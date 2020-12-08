package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

func KeyVaultChildID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	if _, err := azure.ParseKeyVaultChildID(v); err != nil {
		errors = append(errors, fmt.Errorf("can not parse %q as a Key Vault Child resource id: %v", k, err))
	}

	return warnings, errors
}
