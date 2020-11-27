package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

// NOTE: there's some nice things we can do with this around validation
// since these top level objects exist

type RouteFilterResourceID struct {
	ResourceGroup string
	Name          string
}

func ParseRouteFilterID(input string) (*RouteFilterResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Route Filter ID %q: %+v", input, err)
	}

	routeFilter := RouteFilterResourceID{
		ResourceGroup: id.ResourceGroup,
	}

	if routeFilter.Name, err = id.PopSegment("routeFilters"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &routeFilter, nil
}
