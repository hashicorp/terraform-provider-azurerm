package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ReverseReplicationProviderSpecificInput = A2AReprotectInput{}

type A2AReprotectInput struct {
	PolicyId                  *string                  `json:"policyId,omitempty"`
	RecoveryAvailabilitySetId *string                  `json:"recoveryAvailabilitySetId,omitempty"`
	RecoveryCloudServiceId    *string                  `json:"recoveryCloudServiceId,omitempty"`
	RecoveryContainerId       *string                  `json:"recoveryContainerId,omitempty"`
	RecoveryResourceGroupId   *string                  `json:"recoveryResourceGroupId,omitempty"`
	VMDisks                   *[]A2AVMDiskInputDetails `json:"vmDisks,omitempty"`

	// Fields inherited from ReverseReplicationProviderSpecificInput
}

var _ json.Marshaler = A2AReprotectInput{}

func (s A2AReprotectInput) MarshalJSON() ([]byte, error) {
	type wrapper A2AReprotectInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling A2AReprotectInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling A2AReprotectInput: %+v", err)
	}
	decoded["instanceType"] = "A2A"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling A2AReprotectInput: %+v", err)
	}

	return encoded, nil
}
