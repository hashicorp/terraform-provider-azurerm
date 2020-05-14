package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SlotVirtualNetworkSwitchConnectionId struct {
	VirtualNetworkSwitchConnectionId
	SlotName string
}

func SlotVirtualNetworkSwitchConnectionID(resourceId string) (*SlotVirtualNetworkSwitchConnectionId, error) {
	id, err := azure.ParseAzureResourceID(resourceId)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Azure Resource ID %q", id)
	}

	virtualNetworkId, err := VirtualNetworkSwitchConnectionID(resourceId)
	if err != nil {
		return nil, err
	}

	slotVirtualNetworkId := &SlotVirtualNetworkSwitchConnectionId{
		VirtualNetworkSwitchConnectionId: *virtualNetworkId,
	}

	if slotVirtualNetworkId.SlotName, err = id.PopSegment("slots"); err != nil {
		return nil, err
	}

	return slotVirtualNetworkId, nil
}
