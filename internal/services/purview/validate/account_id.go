package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/purview/sdk/2020-12-01-preview/account"
)

func AccountID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := account.ParseAccountID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
