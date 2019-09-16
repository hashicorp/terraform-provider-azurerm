package network

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

// NOTE: there's some nice things we can do with this around validation
// since these top level objects exist

type RouteTableResourceID struct {
	Base azure.ResourceID

	Name string
}

func ParseRouteTableResourceID(input string) (*RouteTableResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Route Table ID %q: %+v", input, err)
	}

	routeTable := RouteTableResourceID{
		Base: *id,
		Name: id.Path["routeTables"],
	}

	if routeTable.Name == "" {
		return nil, fmt.Errorf("ID was missing the `routeTables` element")
	}

	return &routeTable, nil
}
