package backupprotectableitems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ WorkloadProtectableItem = AzureVMWorkloadSAPHanaHSR{}

type AzureVMWorkloadSAPHanaHSR struct {
	IsAutoProtectable       *bool                `json:"isAutoProtectable,omitempty"`
	IsAutoProtected         *bool                `json:"isAutoProtected,omitempty"`
	ParentName              *string              `json:"parentName,omitempty"`
	ParentUniqueName        *string              `json:"parentUniqueName,omitempty"`
	Prebackupvalidation     *PreBackupValidation `json:"prebackupvalidation,omitempty"`
	ServerName              *string              `json:"serverName,omitempty"`
	Subinquireditemcount    *int64               `json:"subinquireditemcount,omitempty"`
	Subprotectableitemcount *int64               `json:"subprotectableitemcount,omitempty"`

	// Fields inherited from WorkloadProtectableItem
	BackupManagementType *string           `json:"backupManagementType,omitempty"`
	FriendlyName         *string           `json:"friendlyName,omitempty"`
	ProtectionState      *ProtectionStatus `json:"protectionState,omitempty"`
	WorkloadType         *string           `json:"workloadType,omitempty"`
}

var _ json.Marshaler = AzureVMWorkloadSAPHanaHSR{}

func (s AzureVMWorkloadSAPHanaHSR) MarshalJSON() ([]byte, error) {
	type wrapper AzureVMWorkloadSAPHanaHSR
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureVMWorkloadSAPHanaHSR: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureVMWorkloadSAPHanaHSR: %+v", err)
	}
	decoded["protectableItemType"] = "SAPHanaHSR"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureVMWorkloadSAPHanaHSR: %+v", err)
	}

	return encoded, nil
}
