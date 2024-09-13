package replicationpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ PolicyProviderSpecificInput = VMwareCbtPolicyCreationInput{}

type VMwareCbtPolicyCreationInput struct {
	AppConsistentFrequencyInMinutes   *int64 `json:"appConsistentFrequencyInMinutes,omitempty"`
	CrashConsistentFrequencyInMinutes *int64 `json:"crashConsistentFrequencyInMinutes,omitempty"`
	RecoveryPointHistoryInMinutes     *int64 `json:"recoveryPointHistoryInMinutes,omitempty"`

	// Fields inherited from PolicyProviderSpecificInput
}

var _ json.Marshaler = VMwareCbtPolicyCreationInput{}

func (s VMwareCbtPolicyCreationInput) MarshalJSON() ([]byte, error) {
	type wrapper VMwareCbtPolicyCreationInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling VMwareCbtPolicyCreationInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling VMwareCbtPolicyCreationInput: %+v", err)
	}
	decoded["instanceType"] = "VMwareCbt"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling VMwareCbtPolicyCreationInput: %+v", err)
	}

	return encoded, nil
}
