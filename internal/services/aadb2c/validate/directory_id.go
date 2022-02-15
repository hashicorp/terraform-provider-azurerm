package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/aadb2c/sdk/2021-04-01-preview/tenants"
)

func B2CDirectoryID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := tenants.ParseB2CDirectoryID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
