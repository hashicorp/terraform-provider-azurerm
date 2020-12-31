package validate

import (
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

func SubscriptionID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}
	id, err := azure.ParseAzureResourceID(v)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q expected to be valid subscription ID, got %q", k, v))
		return
	}

	if _, err := uuid.ParseUUID(id.SubscriptionID); err != nil {
		errors = append(errors, fmt.Errorf("expected subscription id value in %q to be a valid UUID, got %v", v, id.SubscriptionID))
	}

	if id.ResourceGroup != "" || len(id.Path) > 0 {
		errors = append(errors, fmt.Errorf("%q expected to be valid subscription ID, got other ID type %q", k, v))
		return
	}

	return
}
