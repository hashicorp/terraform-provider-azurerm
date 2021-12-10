package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/maps/sdk/2021-02-01/accounts"
)

func AccountID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := accounts.ParseAccountID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
