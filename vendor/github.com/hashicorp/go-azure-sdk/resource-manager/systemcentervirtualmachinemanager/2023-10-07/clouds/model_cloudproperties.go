package clouds

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CloudProperties struct {
	CloudCapacity      *CloudCapacity      `json:"cloudCapacity,omitempty"`
	CloudName          *string             `json:"cloudName,omitempty"`
	InventoryItemId    *string             `json:"inventoryItemId,omitempty"`
	ProvisioningState  *ProvisioningState  `json:"provisioningState,omitempty"`
	StorageQoSPolicies *[]StorageQosPolicy `json:"storageQoSPolicies,omitempty"`
	Uuid               *string             `json:"uuid,omitempty"`
	VMmServerId        *string             `json:"vmmServerId,omitempty"`
}
