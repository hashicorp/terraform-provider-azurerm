package validate

import (
	"fmt"

	dataLakeParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datalake/parse"
)

func AccountID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	if _, err := dataLakeParse.AccountID(v); err != nil {
		errors = append(errors, fmt.Errorf("can not parse %q as a Data Lake Store id: %v", k, err))
	}

	return warnings, errors
}
