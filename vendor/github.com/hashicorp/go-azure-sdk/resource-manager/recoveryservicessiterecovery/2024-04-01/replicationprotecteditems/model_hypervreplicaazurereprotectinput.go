package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ReverseReplicationProviderSpecificInput = HyperVReplicaAzureReprotectInput{}

type HyperVReplicaAzureReprotectInput struct {
	HvHostVMId          *string `json:"hvHostVmId,omitempty"`
	LogStorageAccountId *string `json:"logStorageAccountId,omitempty"`
	OsType              *string `json:"osType,omitempty"`
	StorageAccountId    *string `json:"storageAccountId,omitempty"`
	VHDId               *string `json:"vHDId,omitempty"`
	VirtualMachineName  *string `json:"vmName,omitempty"`

	// Fields inherited from ReverseReplicationProviderSpecificInput

	InstanceType string `json:"instanceType"`
}

func (s HyperVReplicaAzureReprotectInput) ReverseReplicationProviderSpecificInput() BaseReverseReplicationProviderSpecificInputImpl {
	return BaseReverseReplicationProviderSpecificInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = HyperVReplicaAzureReprotectInput{}

func (s HyperVReplicaAzureReprotectInput) MarshalJSON() ([]byte, error) {
	type wrapper HyperVReplicaAzureReprotectInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling HyperVReplicaAzureReprotectInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling HyperVReplicaAzureReprotectInput: %+v", err)
	}

	decoded["instanceType"] = "HyperVReplicaAzure"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling HyperVReplicaAzureReprotectInput: %+v", err)
	}

	return encoded, nil
}
