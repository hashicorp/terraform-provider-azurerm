package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type PoolId struct {
	ResourceGroup    string
	BatchAccountName string
	PoolName         string
}

func PoolID(input string) (*PoolId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Batch Pool ID %q: %+v", input, err)
	}

	pool := PoolId{
		ResourceGroup: id.ResourceGroup,
	}

	if pool.BatchAccountName, err = id.PopSegment("batchAccounts"); err != nil {
		return nil, err
	}

	if pool.PoolName, err = id.PopSegment("pools"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &pool, nil
}
