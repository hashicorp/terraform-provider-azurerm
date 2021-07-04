package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ExpressRouteCircuitConnectionId struct {
	SubscriptionId          string
	ResourceGroup           string
	ExpressRouteCircuitName string
	PeeringName             string
	ConnectionName          string
}

func NewExpressRouteCircuitConnectionID(subscriptionId, resourceGroup, expressRouteCircuitName, peeringName, connectionName string) ExpressRouteCircuitConnectionId {
	return ExpressRouteCircuitConnectionId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		ExpressRouteCircuitName: expressRouteCircuitName,
		PeeringName:             peeringName,
		ConnectionName:          connectionName,
	}
}

func (id ExpressRouteCircuitConnectionId) String() string {
	segments := []string{
		fmt.Sprintf("Connection Name %q", id.ConnectionName),
		fmt.Sprintf("Peering Name %q", id.PeeringName),
		fmt.Sprintf("Express Route Circuit Name %q", id.ExpressRouteCircuitName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Express Route Circuit Connection", segmentsStr)
}

func (id ExpressRouteCircuitConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRouteCircuits/%s/peerings/%s/connections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ExpressRouteCircuitName, id.PeeringName, id.ConnectionName)
}

// ExpressRouteCircuitConnectionID parses a ExpressRouteCircuitConnection ID into an ExpressRouteCircuitConnectionId struct
func ExpressRouteCircuitConnectionID(input string) (*ExpressRouteCircuitConnectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ExpressRouteCircuitConnectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ExpressRouteCircuitName, err = id.PopSegment("expressRouteCircuits"); err != nil {
		return nil, err
	}
	if resourceId.PeeringName, err = id.PopSegment("peerings"); err != nil {
		return nil, err
	}
	if resourceId.ConnectionName, err = id.PopSegment("connections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
