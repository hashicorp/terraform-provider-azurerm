package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/sdk/2022-04-01/resourceguards"
)

func ResourceGuardID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := resourceguards.ParseResourceGuardID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
