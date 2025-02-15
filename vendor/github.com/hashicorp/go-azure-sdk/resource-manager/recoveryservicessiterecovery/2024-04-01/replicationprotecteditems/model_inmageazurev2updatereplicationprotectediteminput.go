package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ UpdateReplicationProtectedItemProviderInput = InMageAzureV2UpdateReplicationProtectedItemInput{}

type InMageAzureV2UpdateReplicationProtectedItemInput struct {
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

	InstanceType string `json:"instanceType"`
}

func (s InMageAzureV2UpdateReplicationProtectedItemInput) UpdateReplicationProtectedItemProviderInput() BaseUpdateReplicationProtectedItemProviderInputImpl {
	return BaseUpdateReplicationProtectedItemProviderInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = InMageAzureV2UpdateReplicationProtectedItemInput{}

func (s InMageAzureV2UpdateReplicationProtectedItemInput) MarshalJSON() ([]byte, error) {
	type wrapper InMageAzureV2UpdateReplicationProtectedItemInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageAzureV2UpdateReplicationProtectedItemInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageAzureV2UpdateReplicationProtectedItemInput: %+v", err)
	}

	decoded["instanceType"] = "InMageAzureV2"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageAzureV2UpdateReplicationProtectedItemInput: %+v", err)
	}

	return encoded, nil
}
