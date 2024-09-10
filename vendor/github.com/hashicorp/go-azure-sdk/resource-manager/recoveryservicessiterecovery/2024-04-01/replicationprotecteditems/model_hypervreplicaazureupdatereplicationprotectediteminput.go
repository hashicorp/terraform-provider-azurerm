package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ UpdateReplicationProtectedItemProviderInput = HyperVReplicaAzureUpdateReplicationProtectedItemInput{}

type HyperVReplicaAzureUpdateReplicationProtectedItemInput struct {
	DiskIdToDiskEncryptionMap       *map[string]string    `json:"diskIdToDiskEncryptionMap,omitempty"`
	RecoveryAzureV1ResourceGroupId  *string               `json:"recoveryAzureV1ResourceGroupId,omitempty"`
	RecoveryAzureV2ResourceGroupId  *string               `json:"recoveryAzureV2ResourceGroupId,omitempty"`
	SqlServerLicenseType            *SqlServerLicenseType `json:"sqlServerLicenseType,omitempty"`
	TargetAvailabilityZone          *string               `json:"targetAvailabilityZone,omitempty"`
	TargetManagedDiskTags           *map[string]string    `json:"targetManagedDiskTags,omitempty"`
	TargetNicTags                   *map[string]string    `json:"targetNicTags,omitempty"`
	TargetProximityPlacementGroupId *string               `json:"targetProximityPlacementGroupId,omitempty"`
	TargetVMTags                    *map[string]string    `json:"targetVmTags,omitempty"`
	UseManagedDisks                 *string               `json:"useManagedDisks,omitempty"`
	VMDisks                         *[]UpdateDiskInput    `json:"vmDisks,omitempty"`

	// Fields inherited from UpdateReplicationProtectedItemProviderInput
}

var _ json.Marshaler = HyperVReplicaAzureUpdateReplicationProtectedItemInput{}

func (s HyperVReplicaAzureUpdateReplicationProtectedItemInput) MarshalJSON() ([]byte, error) {
	type wrapper HyperVReplicaAzureUpdateReplicationProtectedItemInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling HyperVReplicaAzureUpdateReplicationProtectedItemInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling HyperVReplicaAzureUpdateReplicationProtectedItemInput: %+v", err)
	}
	decoded["instanceType"] = "HyperVReplicaAzure"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling HyperVReplicaAzureUpdateReplicationProtectedItemInput: %+v", err)
	}

	return encoded, nil
}
