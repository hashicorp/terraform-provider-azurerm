package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type RouteTableId struct {
	ResourceGroup string
	Name          string
}

func RouteTableID(input string) (*RouteTableId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Route Table ID %q: %+v", input, err)
	}

	routeTable := RouteTableId{
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
