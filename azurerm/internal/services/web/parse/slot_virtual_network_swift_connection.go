package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SlotVirtualNetworkSwiftConnectionId struct {
	VirtualNetworkSwiftConnectionId
	SlotName string
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
		VirtualNetworkSwiftConnectionId: *virtualNetworkId,
	}

	if slotVirtualNetworkId.SlotName, err = id.PopSegment("slots"); err != nil {
		return nil, err
	}

	return slotVirtualNetworkId, nil
}
