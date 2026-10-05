package protectionpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ProtectionPolicy = AzureVMWorkloadProtectionPolicy{}

type AzureVMWorkloadProtectionPolicy struct {
	MakePolicyConsistent *bool                  `json:"makePolicyConsistent,omitempty"`
	Settings             *Settings              `json:"settings,omitempty"`
	SubProtectionPolicy  *[]SubProtectionPolicy `json:"subProtectionPolicy,omitempty"`
	WorkLoadType         *WorkloadType          `json:"workLoadType,omitempty"`

	// Fields inherited from ProtectionPolicy

	BackupManagementType           string    `json:"backupManagementType"`
	ProtectedItemsCount            *int64    `json:"protectedItemsCount,omitempty"`
	ResourceGuardOperationRequests *[]string `json:"resourceGuardOperationRequests,omitempty"`
}

func (s AzureVMWorkloadProtectionPolicy) ProtectionPolicy() BaseProtectionPolicyImpl {
	return BaseProtectionPolicyImpl{
		BackupManagementType:           s.BackupManagementType,
		ProtectedItemsCount:            s.ProtectedItemsCount,
		ResourceGuardOperationRequests: s.ResourceGuardOperationRequests,
	}
}

var _ json.Marshaler = AzureVMWorkloadProtectionPolicy{}

func (s AzureVMWorkloadProtectionPolicy) MarshalJSON() ([]byte, error) {
	type wrapper AzureVMWorkloadProtectionPolicy
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureVMWorkloadProtectionPolicy: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureVMWorkloadProtectionPolicy: %+v", err)
	}

	decoded["backupManagementType"] = "AzureWorkload"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureVMWorkloadProtectionPolicy: %+v", err)
	}

	return encoded, nil
}
