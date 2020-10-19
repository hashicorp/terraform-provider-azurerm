package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualHubRouteTableId struct {
	ResourceGroup  string
	VirtualHubName string
	Name           string
}

func VirtualHubRouteTableID(input string) (*VirtualHubRouteTableId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing virtualHubRouteTable ID %q: %+v", input, err)
	}

	virtualHubRouteTable := VirtualHubRouteTableId{
		ResourceGroup: id.ResourceGroup,
	}

	if virtualHubRouteTable.VirtualHubName, err = id.PopSegment("virtualHubs"); err != nil {
		return nil, err
	}

	if virtualHubRouteTable.Name, err = id.PopSegment("hubRouteTables"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &virtualHubRouteTable, nil
}
