package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type PeeringId struct {
	SubscriptionId          string
	ResourceGroup           string
	ExpressRouteCircuitName string
	Name                    string
}

func NewPeeringID(subscriptionId, resourceGroup, expressRouteCircuitName, name string) PeeringId {
	return PeeringId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		ExpressRouteCircuitName: expressRouteCircuitName,
		Name:                    name,
	}
}

func (id PeeringId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Express Route Circuit Name %q", id.ExpressRouteCircuitName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Peering", segmentsStr)
}

func (id PeeringId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRouteCircuits/%s/peerings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ExpressRouteCircuitName, id.Name)
}

// PeeringID parses a Peering ID into an PeeringId struct
func PeeringID(input string) (*PeeringId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := PeeringId{
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
	if resourceId.Name, err = id.PopSegment("peerings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
