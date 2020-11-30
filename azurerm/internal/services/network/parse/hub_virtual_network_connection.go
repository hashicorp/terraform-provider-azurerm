package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type HubVirtualNetworkConnectionId struct {
	SubscriptionId string
	ResourceGroup  string
	VirtualHubName string
	Name           string
}

func NewHubVirtualNetworkConnectionID(subscriptionId, resourceGroup, virtualHubName, name string) HubVirtualNetworkConnectionId {
	return HubVirtualNetworkConnectionId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		VirtualHubName: virtualHubName,
		Name:           name,
	}
}

func (id HubVirtualNetworkConnectionId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualHubs/%s/hubVirtualNetworkConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VirtualHubName, id.Name)
}

// HubVirtualNetworkConnectionID parses a HubVirtualNetworkConnection ID into an HubVirtualNetworkConnectionId struct
func HubVirtualNetworkConnectionID(input string) (*HubVirtualNetworkConnectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := HubVirtualNetworkConnectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.VirtualHubName, err = id.PopSegment("virtualHubs"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("hubVirtualNetworkConnections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
