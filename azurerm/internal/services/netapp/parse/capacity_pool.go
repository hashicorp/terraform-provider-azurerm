package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type CapacityPoolId struct {
	ResourceGroup     string
	NetAppAccountName string
	Name              string
}

func CapacityPoolID(input string) (*CapacityPoolId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse NetApp Pool ID %q: %+v", input, err)
	}

	service := CapacityPoolId{
		ResourceGroup: id.ResourceGroup,
	}

	if service.NetAppAccountName, err = id.PopSegment("netAppAccounts"); err != nil {
		return nil, err
	}

	if service.Name, err = id.PopSegment("capacityPools"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &service, nil
}
