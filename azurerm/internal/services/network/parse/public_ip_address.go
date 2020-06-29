package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type PublicIPAddressId struct {
	ResourceGroup string
	Name          string
}

func PublicIPAddressID(input string) (*PublicIPAddressId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Public IP Address ID %q: %+v", input, err)
	}

	ipAddress := PublicIPAddressId{
		ResourceGroup: id.ResourceGroup,
	}

	if ipAddress.Name, err = id.PopSegment("publicIPAddresses"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &ipAddress, nil
}
