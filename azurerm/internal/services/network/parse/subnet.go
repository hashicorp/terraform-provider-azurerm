package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SubnetId struct {
	ResourceGroup      string
	VirtualNetworkName string
	Name               string
}

func SubnetID(input string) (*SubnetId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Subnet ID %q: %+v", input, err)
	}

	subnet := SubnetId{
		ResourceGroup: id.ResourceGroup,
	}

	if subnet.VirtualNetworkName, err = id.PopSegment("virtualNetworks"); err != nil {
		return nil, err
	}

	if subnet.Name, err = id.PopSegment("subnets"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &subnet, nil
}
