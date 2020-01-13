package network

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

// NOTE: there's some nice things we can do with this around validation
// since these top level objects exist

type RouteTableResourceID struct {
	ResourceGroup string
	Name          string
}

func ParseRouteTableID(input string) (*RouteTableResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Route Table ID %q: %+v", input, err)
	}

	routeTable := RouteTableResourceID{
		ResourceGroup: id.ResourceGroup,
	}

	if routeTable.Name, err = id.PopSegment("routeTables"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &routeTable, nil
}
