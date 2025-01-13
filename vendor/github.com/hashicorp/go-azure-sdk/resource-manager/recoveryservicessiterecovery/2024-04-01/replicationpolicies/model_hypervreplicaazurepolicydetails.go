package replicationpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ PolicyProviderSpecificDetails = HyperVReplicaAzurePolicyDetails{}

type HyperVReplicaAzurePolicyDetails struct {
	ActiveStorageAccountId                        *string `json:"activeStorageAccountId,omitempty"`
	ApplicationConsistentSnapshotFrequencyInHours *int64  `json:"applicationConsistentSnapshotFrequencyInHours,omitempty"`
	Encryption                                    *string `json:"encryption,omitempty"`
	OnlineReplicationStartTime                    *string `json:"onlineReplicationStartTime,omitempty"`
	RecoveryPointHistoryDurationInHours           *int64  `json:"recoveryPointHistoryDurationInHours,omitempty"`
	ReplicationInterval                           *int64  `json:"replicationInterval,omitempty"`

	// Fields inherited from PolicyProviderSpecificDetails

	InstanceType string `json:"instanceType"`
}

func (s HyperVReplicaAzurePolicyDetails) PolicyProviderSpecificDetails() BasePolicyProviderSpecificDetailsImpl {
	return BasePolicyProviderSpecificDetailsImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = HyperVReplicaAzurePolicyDetails{}

func (s HyperVReplicaAzurePolicyDetails) MarshalJSON() ([]byte, error) {
	type wrapper HyperVReplicaAzurePolicyDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling HyperVReplicaAzurePolicyDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling HyperVReplicaAzurePolicyDetails: %+v", err)
	}

	decoded["instanceType"] = "HyperVReplicaAzure"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling HyperVReplicaAzurePolicyDetails: %+v", err)
	}

	return encoded, nil
}
