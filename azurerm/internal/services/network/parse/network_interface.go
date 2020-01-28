package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type NetworkInterfaceId struct {
	ResourceGroup string
	Name          string
}

func NetworkInterfaceID(input string) (*NetworkInterfaceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Network Interface ID %q: %+v", input, err)
	}

	server := NetworkInterfaceId{
		ResourceGroup: id.ResourceGroup,
	}

	if server.Name, err = id.PopSegment("networkInterfaces"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &server, nil
}
