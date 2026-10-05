package virtualnetworks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkProperties struct {
	InventoryItemId   *string            `json:"inventoryItemId,omitempty"`
	NetworkName       *string            `json:"networkName,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	Uuid              *string            `json:"uuid,omitempty"`
	VMmServerId       *string            `json:"vmmServerId,omitempty"`
}
