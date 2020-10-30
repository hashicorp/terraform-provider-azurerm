package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ExpressRouteConnectionId struct {
	ResourceGroup           string
	ExpressRouteGatewayName string
	Name                    string
}

func NewExpressRouteConnectionID(resourceGroup string, expressRouteGatewayName string, name string) ExpressRouteConnectionId {
	return ExpressRouteConnectionId{
		ResourceGroup:           resourceGroup,
		ExpressRouteGatewayName: expressRouteGatewayName,
		Name:                    name,
	}
}

func (id ExpressRouteConnectionId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRouteGateways/%s/expressRouteConnections/%s", subscriptionId, id.ResourceGroup, id.ExpressRouteGatewayName, id.Name)
}

func ExpressRouteConnectionID(input string) (*ExpressRouteConnectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing expressRouteConnection ID %q: %+v", input, err)
	}

	expressRouteConnection := ExpressRouteConnectionId{
		ResourceGroup: id.ResourceGroup,
	}

	if expressRouteConnection.ExpressRouteGatewayName, err = id.PopSegment("expressRouteGateways"); err != nil {
		return nil, err
	}

	if expressRouteConnection.Name, err = id.PopSegment("expressRouteConnections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &expressRouteConnection, nil
}
