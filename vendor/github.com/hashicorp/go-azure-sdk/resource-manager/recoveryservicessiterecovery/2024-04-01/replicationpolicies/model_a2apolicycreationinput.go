package replicationpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ PolicyProviderSpecificInput = A2APolicyCreationInput{}

type A2APolicyCreationInput struct {
	AppConsistentFrequencyInMinutes   *int64               `json:"appConsistentFrequencyInMinutes,omitempty"`
	CrashConsistentFrequencyInMinutes *int64               `json:"crashConsistentFrequencyInMinutes,omitempty"`
	MultiVMSyncStatus                 SetMultiVMSyncStatus `json:"multiVmSyncStatus"`
	RecoveryPointHistory              *int64               `json:"recoveryPointHistory,omitempty"`

	// Fields inherited from PolicyProviderSpecificInput

	InstanceType string `json:"instanceType"`
}

func (s A2APolicyCreationInput) PolicyProviderSpecificInput() BasePolicyProviderSpecificInputImpl {
	return BasePolicyProviderSpecificInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = A2APolicyCreationInput{}

func (s A2APolicyCreationInput) MarshalJSON() ([]byte, error) {
	type wrapper A2APolicyCreationInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling A2APolicyCreationInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling A2APolicyCreationInput: %+v", err)
	}

	decoded["instanceType"] = "A2A"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling A2APolicyCreationInput: %+v", err)
	}

	return encoded, nil
}
