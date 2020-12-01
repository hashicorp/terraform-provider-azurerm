package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualNetworkSwiftConnectionId struct {
	SubscriptionId    string
	ResourceGroup     string
	SiteName          string
	NetworkConfigName string
}

func NewVirtualNetworkSwiftConnectionID(subscriptionId, resourceGroup, siteName, networkConfigName string) VirtualNetworkSwiftConnectionId {
	return VirtualNetworkSwiftConnectionId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		SiteName:          siteName,
		NetworkConfigName: networkConfigName,
	}
}

func (id VirtualNetworkSwiftConnectionId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/networkConfig/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SiteName, id.NetworkConfigName)
}

// VirtualNetworkSwiftConnectionID parses a VirtualNetworkSwiftConnection ID into an VirtualNetworkSwiftConnectionId struct
func VirtualNetworkSwiftConnectionID(input string) (*VirtualNetworkSwiftConnectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := VirtualNetworkSwiftConnectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SiteName, err = id.PopSegment("sites"); err != nil {
		return nil, err
	}
	if resourceId.NetworkConfigName, err = id.PopSegment("networkConfig"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
