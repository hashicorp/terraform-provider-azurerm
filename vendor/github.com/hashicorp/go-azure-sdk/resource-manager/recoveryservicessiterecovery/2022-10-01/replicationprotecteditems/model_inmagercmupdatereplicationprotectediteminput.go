package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ UpdateReplicationProtectedItemProviderInput = InMageRcmUpdateReplicationProtectedItemInput{}

type InMageRcmUpdateReplicationProtectedItemInput struct {
	LicenseType                           *LicenseType         `json:"licenseType,omitempty"`
	TargetAvailabilitySetId               *string              `json:"targetAvailabilitySetId,omitempty"`
	TargetAvailabilityZone                *string              `json:"targetAvailabilityZone,omitempty"`
	TargetBootDiagnosticsStorageAccountId *string              `json:"targetBootDiagnosticsStorageAccountId,omitempty"`
	TargetNetworkId                       *string              `json:"targetNetworkId,omitempty"`
	TargetProximityPlacementGroupId       *string              `json:"targetProximityPlacementGroupId,omitempty"`
	TargetResourceGroupId                 *string              `json:"targetResourceGroupId,omitempty"`
	TargetVMName                          *string              `json:"targetVmName,omitempty"`
	TargetVMSize                          *string              `json:"targetVmSize,omitempty"`
	TestNetworkId                         *string              `json:"testNetworkId,omitempty"`
	VMNics                                *[]InMageRcmNicInput `json:"vmNics,omitempty"`

	// Fields inherited from UpdateReplicationProtectedItemProviderInput
}

var _ json.Marshaler = InMageRcmUpdateReplicationProtectedItemInput{}

func (s InMageRcmUpdateReplicationProtectedItemInput) MarshalJSON() ([]byte, error) {
	type wrapper InMageRcmUpdateReplicationProtectedItemInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageRcmUpdateReplicationProtectedItemInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageRcmUpdateReplicationProtectedItemInput: %+v", err)
	}
	decoded["instanceType"] = "InMageRcm"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageRcmUpdateReplicationProtectedItemInput: %+v", err)
	}

	return encoded, nil
}
