package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/parse"
)

func RestorableDatabaseID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	if _, err := parse.RestorableDroppedDatabaseID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a MsSql Restorable Database resource id: %v", k, err))
	}

	return warnings, errors
}
