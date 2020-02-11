package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type BatchPoolId struct {
	ResourceGroup string
	AccountName   string
	Name          string
}

func BatchPoolID(input string) (*BatchPoolId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Batch Pool ID %q: %+v", input, err)
	}

	pool := BatchPoolId{
		ResourceGroup: id.ResourceGroup,
	}

	if pool.AccountName, err = id.PopSegment("batchAccounts"); err != nil {
		return nil, err
	}

	if pool.Name, err = id.PopSegment("pools"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &pool, nil
}
