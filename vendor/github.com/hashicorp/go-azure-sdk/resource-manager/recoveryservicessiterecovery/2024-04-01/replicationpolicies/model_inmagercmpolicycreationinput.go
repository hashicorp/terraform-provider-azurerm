package replicationpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ PolicyProviderSpecificInput = InMageRcmPolicyCreationInput{}

type InMageRcmPolicyCreationInput struct {
	AppConsistentFrequencyInMinutes   *int64  `json:"appConsistentFrequencyInMinutes,omitempty"`
	CrashConsistentFrequencyInMinutes *int64  `json:"crashConsistentFrequencyInMinutes,omitempty"`
	EnableMultiVMSync                 *string `json:"enableMultiVmSync,omitempty"`
	RecoveryPointHistoryInMinutes     *int64  `json:"recoveryPointHistoryInMinutes,omitempty"`

	// Fields inherited from PolicyProviderSpecificInput

	InstanceType string `json:"instanceType"`
}

func (s InMageRcmPolicyCreationInput) PolicyProviderSpecificInput() BasePolicyProviderSpecificInputImpl {
	return BasePolicyProviderSpecificInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = InMageRcmPolicyCreationInput{}

func (s InMageRcmPolicyCreationInput) MarshalJSON() ([]byte, error) {
	type wrapper InMageRcmPolicyCreationInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageRcmPolicyCreationInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageRcmPolicyCreationInput: %+v", err)
	}

	decoded["instanceType"] = "InMageRcm"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageRcmPolicyCreationInput: %+v", err)
	}

	return encoded, nil
}
