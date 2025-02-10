package backupprotectableitems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ WorkloadProtectableItem = IaaSVMProtectableItem{}

type IaaSVMProtectableItem struct {
	ResourceGroup         *string `json:"resourceGroup,omitempty"`
	VirtualMachineId      *string `json:"virtualMachineId,omitempty"`
	VirtualMachineVersion *string `json:"virtualMachineVersion,omitempty"`

	// Fields inherited from WorkloadProtectableItem

	BackupManagementType *string           `json:"backupManagementType,omitempty"`
	FriendlyName         *string           `json:"friendlyName,omitempty"`
	ProtectableItemType  string            `json:"protectableItemType"`
	ProtectionState      *ProtectionStatus `json:"protectionState,omitempty"`
	WorkloadType         *string           `json:"workloadType,omitempty"`
}

func (s IaaSVMProtectableItem) WorkloadProtectableItem() BaseWorkloadProtectableItemImpl {
	return BaseWorkloadProtectableItemImpl{
		BackupManagementType: s.BackupManagementType,
		FriendlyName:         s.FriendlyName,
		ProtectableItemType:  s.ProtectableItemType,
		ProtectionState:      s.ProtectionState,
		WorkloadType:         s.WorkloadType,
	}
}

var _ json.Marshaler = IaaSVMProtectableItem{}

func (s IaaSVMProtectableItem) MarshalJSON() ([]byte, error) {
	type wrapper IaaSVMProtectableItem
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling IaaSVMProtectableItem: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling IaaSVMProtectableItem: %+v", err)
	}

	decoded["protectableItemType"] = "IaaSVMProtectableItem"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling IaaSVMProtectableItem: %+v", err)
	}

	return encoded, nil
}
