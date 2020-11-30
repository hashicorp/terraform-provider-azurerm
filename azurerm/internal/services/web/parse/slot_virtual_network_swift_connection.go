package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SlotVirtualNetworkSwiftConnectionId struct {
	SubscriptionId    string
	ResourceGroup     string
	SiteName          string
	SlotName          string
	NetworkconfigName string
}

func NewSlotVirtualNetworkSwiftConnectionID(subscriptionId, resourceGroup, siteName, slotName, networkconfigName string) SlotVirtualNetworkSwiftConnectionId {
	return SlotVirtualNetworkSwiftConnectionId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		SiteName:          siteName,
		SlotName:          slotName,
		NetworkconfigName: networkconfigName,
	}
}

func (id SlotVirtualNetworkSwiftConnectionId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/slots/%s/networkconfig/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SiteName, id.SlotName, id.NetworkconfigName)
}

// SlotVirtualNetworkSwiftConnectionID parses a SlotVirtualNetworkSwiftConnection ID into an SlotVirtualNetworkSwiftConnectionId struct
func SlotVirtualNetworkSwiftConnectionID(input string) (*SlotVirtualNetworkSwiftConnectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SlotVirtualNetworkSwiftConnectionId{
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
	if resourceId.SlotName, err = id.PopSegment("slots"); err != nil {
		return nil, err
	}
	if resourceId.NetworkconfigName, err = id.PopSegment("networkconfig"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
