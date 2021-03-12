package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ExpressRouteConnectionId struct {
	SubscriptionId          string
	ResourceGroup           string
	ExpressRouteGatewayName string
	Name                    string
}

func NewExpressRouteConnectionID(subscriptionId, resourceGroup, expressRouteGatewayName, name string) ExpressRouteConnectionId {
	return ExpressRouteConnectionId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		ExpressRouteGatewayName: expressRouteGatewayName,
		Name:                    name,
	}
}

func (id ExpressRouteConnectionId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Express Route Gateway Name %q", id.ExpressRouteGatewayName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Express Route Connection", segmentsStr)
}

func (id ExpressRouteConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRouteGateways/%s/expressRouteConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ExpressRouteGatewayName, id.Name)
}

// ExpressRouteConnectionID parses a ExpressRouteConnection ID into an ExpressRouteConnectionId struct
func ExpressRouteConnectionID(input string) (*ExpressRouteConnectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ExpressRouteConnectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ExpressRouteGatewayName, err = id.PopSegment("expressRouteGateways"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("expressRouteConnections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
