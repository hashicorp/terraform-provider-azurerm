package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ExpressRouteCircuitPeeringId struct {
	SubscriptionId          string
	ResourceGroup           string
	ExpressRouteCircuitName string
	PeeringName             string
}

func NewExpressRouteCircuitPeeringID(subscriptionId, resourceGroup, expressRouteCircuitName, peeringName string) ExpressRouteCircuitPeeringId {
	return ExpressRouteCircuitPeeringId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		ExpressRouteCircuitName: expressRouteCircuitName,
		PeeringName:             peeringName,
	}
}

func (id ExpressRouteCircuitPeeringId) String() string {
	segments := []string{
		fmt.Sprintf("Peering Name %q", id.PeeringName),
		fmt.Sprintf("Express Route Circuit Name %q", id.ExpressRouteCircuitName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Express Route Circuit Peering", segmentsStr)
}

func (id ExpressRouteCircuitPeeringId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRouteCircuits/%s/peerings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ExpressRouteCircuitName, id.PeeringName)
}

// ExpressRouteCircuitPeeringID parses a ExpressRouteCircuitPeering ID into an ExpressRouteCircuitPeeringId struct
func ExpressRouteCircuitPeeringID(input string) (*ExpressRouteCircuitPeeringId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ExpressRouteCircuitPeeringId{
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

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
