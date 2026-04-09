package replicationpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ PolicyProviderSpecificInput = HyperVReplicaAzurePolicyInput{}

type HyperVReplicaAzurePolicyInput struct {
	ApplicationConsistentSnapshotFrequencyInHours *int64    `json:"applicationConsistentSnapshotFrequencyInHours,omitempty"`
	OnlineReplicationStartTime                    *string   `json:"onlineReplicationStartTime,omitempty"`
	RecoveryPointHistoryDuration                  *int64    `json:"recoveryPointHistoryDuration,omitempty"`
	ReplicationInterval                           *int64    `json:"replicationInterval,omitempty"`
	StorageAccounts                               *[]string `json:"storageAccounts,omitempty"`

	// Fields inherited from PolicyProviderSpecificInput

	InstanceType string `json:"instanceType"`
}

func (s HyperVReplicaAzurePolicyInput) PolicyProviderSpecificInput() BasePolicyProviderSpecificInputImpl {
	return BasePolicyProviderSpecificInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = HyperVReplicaAzurePolicyInput{}

func (s HyperVReplicaAzurePolicyInput) MarshalJSON() ([]byte, error) {
	type wrapper HyperVReplicaAzurePolicyInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling HyperVReplicaAzurePolicyInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling HyperVReplicaAzurePolicyInput: %+v", err)
	}

	decoded["instanceType"] = "HyperVReplicaAzure"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling HyperVReplicaAzurePolicyInput: %+v", err)
	}

	return encoded, nil
}
