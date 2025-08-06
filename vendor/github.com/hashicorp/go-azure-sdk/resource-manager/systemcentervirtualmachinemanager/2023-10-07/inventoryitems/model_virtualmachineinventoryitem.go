package inventoryitems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ InventoryItemProperties = VirtualMachineInventoryItem{}

type VirtualMachineInventoryItem struct {
	BiosGuid                 *string               `json:"biosGuid,omitempty"`
	Cloud                    *InventoryItemDetails `json:"cloud,omitempty"`
	IPAddresses              *[]string             `json:"ipAddresses,omitempty"`
	ManagedMachineResourceId *string               `json:"managedMachineResourceId,omitempty"`
	OsName                   *string               `json:"osName,omitempty"`
	OsType                   *OsType               `json:"osType,omitempty"`
	OsVersion                *string               `json:"osVersion,omitempty"`
	PowerState               *string               `json:"powerState,omitempty"`

	// Fields inherited from InventoryItemProperties

	InventoryItemName *string            `json:"inventoryItemName,omitempty"`
	InventoryType     InventoryType      `json:"inventoryType"`
	ManagedResourceId *string            `json:"managedResourceId,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	Uuid              *string            `json:"uuid,omitempty"`
}

func (s VirtualMachineInventoryItem) InventoryItemProperties() BaseInventoryItemPropertiesImpl {
	return BaseInventoryItemPropertiesImpl{
		InventoryItemName: s.InventoryItemName,
		InventoryType:     s.InventoryType,
		ManagedResourceId: s.ManagedResourceId,
		ProvisioningState: s.ProvisioningState,
		Uuid:              s.Uuid,
	}
}

var _ json.Marshaler = VirtualMachineInventoryItem{}

func (s VirtualMachineInventoryItem) MarshalJSON() ([]byte, error) {
	type wrapper VirtualMachineInventoryItem
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling VirtualMachineInventoryItem: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling VirtualMachineInventoryItem: %+v", err)
	}

	decoded["inventoryType"] = "VirtualMachine"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling VirtualMachineInventoryItem: %+v", err)
	}

	return encoded, nil
}
