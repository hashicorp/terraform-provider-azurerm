package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SlotVirtualNetworkSwiftConnectionId struct {
	ResourceGroup string
	SiteName      string
	SlotName      string
}

func SlotVirtualNetworkSwiftConnectionID(resourceId string) (*SlotVirtualNetworkSwiftConnectionId, error) {
	id, err := azure.ParseAzureResourceID(resourceId)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Azure Resource ID %q", id)
	}

	virtualNetworkId, err := VirtualNetworkSwiftConnectionID(resourceId)
	if err != nil {
		return nil, err
	}

	slotVirtualNetworkId := &SlotVirtualNetworkSwiftConnectionId{
		ResourceGroup: virtualNetworkId.ResourceGroup,
		SiteName:      virtualNetworkId.SiteName,
	}

	if slotVirtualNetworkId.SlotName, err = id.PopSegment("slots"); err != nil {
		return nil, err
	}

	return slotVirtualNetworkId, nil
}
