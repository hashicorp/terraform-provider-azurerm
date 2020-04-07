package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualNetworkSwitchConnectionSlotId struct {
	VirtualNetworkSwitchConnectionId
	SlotName string
}

func VirtualNetworkSwitchConnectionSlotID(ID string) (*VirtualNetworkSwitchConnectionSlotId, error) {
	id, err := azure.ParseAzureResourceID(ID)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Azure Resource ID %q", id)
	}

	virtualNetworkId, err := VirtualNetworkSwitchConnectionID(ID)
	if err != nil {
		return nil, err
	}

	slotVirtualNetworkId := &VirtualNetworkSwitchConnectionSlotId{
		VirtualNetworkSwitchConnectionId: *virtualNetworkId,
	}

	if slotVirtualNetworkId.SlotName, err = id.PopSegment("slots"); err != nil {
		return nil, err
	}

	return slotVirtualNetworkId, nil
}
