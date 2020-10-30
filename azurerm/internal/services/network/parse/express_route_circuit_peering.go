package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ExpressRouteCircuitPeeringId struct {
	ResourceGroup string
	CircuitName   string
	Name          string
}

func ExpressRouteCircuitPeeringID(input string) (*ExpressRouteCircuitPeeringId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing expressRouteCircuitPeering ID %q: %+v", input, err)
	}

	expressRouteCircuitPeering := ExpressRouteCircuitPeeringId{
		ResourceGroup: id.ResourceGroup,
	}
	if expressRouteCircuitPeering.CircuitName, err = id.PopSegment("expressRouteCircuits"); err != nil {
		return nil, err
	}
	if expressRouteCircuitPeering.Name, err = id.PopSegment("peerings"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &expressRouteCircuitPeering, nil
}
