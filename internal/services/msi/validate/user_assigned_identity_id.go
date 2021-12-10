package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/msi/sdk/2018-11-30/managedidentity"
)

func UserAssignedIdentityID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := managedidentity.ParseUserAssignedIdentitiesID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
