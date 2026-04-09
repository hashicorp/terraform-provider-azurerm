package inventoryitems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ InventoryItemProperties = VirtualMachineTemplateInventoryItem{}

type VirtualMachineTemplateInventoryItem struct {
	CpuCount *int64  `json:"cpuCount,omitempty"`
	MemoryMB *int64  `json:"memoryMB,omitempty"`
	OsName   *string `json:"osName,omitempty"`
	OsType   *OsType `json:"osType,omitempty"`

	// Fields inherited from InventoryItemProperties

	InventoryItemName *string            `json:"inventoryItemName,omitempty"`
	InventoryType     InventoryType      `json:"inventoryType"`
	ManagedResourceId *string            `json:"managedResourceId,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	Uuid              *string            `json:"uuid,omitempty"`
}

func (s VirtualMachineTemplateInventoryItem) InventoryItemProperties() BaseInventoryItemPropertiesImpl {
	return BaseInventoryItemPropertiesImpl{
		InventoryItemName: s.InventoryItemName,
		InventoryType:     s.InventoryType,
		ManagedResourceId: s.ManagedResourceId,
		ProvisioningState: s.ProvisioningState,
		Uuid:              s.Uuid,
	}
}

var _ json.Marshaler = VirtualMachineTemplateInventoryItem{}

func (s VirtualMachineTemplateInventoryItem) MarshalJSON() ([]byte, error) {
	type wrapper VirtualMachineTemplateInventoryItem
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling VirtualMachineTemplateInventoryItem: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling VirtualMachineTemplateInventoryItem: %+v", err)
	}

	decoded["inventoryType"] = "VirtualMachineTemplate"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling VirtualMachineTemplateInventoryItem: %+v", err)
	}

	return encoded, nil
}
