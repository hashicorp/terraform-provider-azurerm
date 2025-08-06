package replicationprotectioncontainers

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SwitchClusterProtectionProviderSpecificInput = A2ASwitchClusterProtectionInput{}

type A2ASwitchClusterProtectionInput struct {
	PolicyId             *string                   `json:"policyId,omitempty"`
	ProtectedItemsDetail *[]A2AProtectedItemDetail `json:"protectedItemsDetail,omitempty"`
	RecoveryContainerId  *string                   `json:"recoveryContainerId,omitempty"`

	// Fields inherited from SwitchClusterProtectionProviderSpecificInput

	InstanceType string `json:"instanceType"`
}

func (s A2ASwitchClusterProtectionInput) SwitchClusterProtectionProviderSpecificInput() BaseSwitchClusterProtectionProviderSpecificInputImpl {
	return BaseSwitchClusterProtectionProviderSpecificInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = A2ASwitchClusterProtectionInput{}

func (s A2ASwitchClusterProtectionInput) MarshalJSON() ([]byte, error) {
	type wrapper A2ASwitchClusterProtectionInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling A2ASwitchClusterProtectionInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling A2ASwitchClusterProtectionInput: %+v", err)
	}

	decoded["instanceType"] = "A2A"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling A2ASwitchClusterProtectionInput: %+v", err)
	}

	return encoded, nil
}
