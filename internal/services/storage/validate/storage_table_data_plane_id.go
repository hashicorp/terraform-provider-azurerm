package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
)

func StorageTableDataPlaneID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	parsed, err := parse.StorageTableDataPlaneID(v)
	if err != nil {
		errors = append(errors, err)
		return
	}

	if _, err := StorageAccountName(parsed.AccountName, key); err != nil {
		errors = append(errors, err...)
	}

	if _, err := StorageTableName(parsed.Name, key); err != nil {
		errors = append(errors, err...)
	}

	return
}
