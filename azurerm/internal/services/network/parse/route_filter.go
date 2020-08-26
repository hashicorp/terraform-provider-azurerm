package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type RouteFilterId struct {
	ResourceGroup string
	Name          string
}

func RouteFilterID(input string) (*RouteFilterId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Route Filter ID %q: %+v", input, err)
	}

	routeFilter := RouteFilterId{
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
