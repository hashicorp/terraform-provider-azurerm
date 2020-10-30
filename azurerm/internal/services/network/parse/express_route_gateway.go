package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ExpressRouteGatewayId struct {
	ResourceGroup string
	Name          string
}

func NewExpressRouteGatewayID(resourceGroup string, name string) ExpressRouteGatewayId {
	return ExpressRouteGatewayId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}

func (id ExpressRouteGatewayId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRouteGateways/%s", subscriptionId, id.ResourceGroup, id.Name)
}

func ExpressRouteGatewayID(input string) (*ExpressRouteGatewayId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing networkExpressRouteGateway ID %q: %+v", input, err)
	}

	expressRouteGateway := ExpressRouteGatewayId{
		ResourceGroup: id.ResourceGroup,
	}

	if expressRouteGateway.Name, err = id.PopSegment("expressRouteGateways"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &expressRouteGateway, nil
}
