package replicationpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ PolicyProviderSpecificInput = InMageRcmFailbackPolicyCreationInput{}

type InMageRcmFailbackPolicyCreationInput struct {
	AppConsistentFrequencyInMinutes   *int64 `json:"appConsistentFrequencyInMinutes,omitempty"`
	CrashConsistentFrequencyInMinutes *int64 `json:"crashConsistentFrequencyInMinutes,omitempty"`

	// Fields inherited from PolicyProviderSpecificInput
}

var _ json.Marshaler = InMageRcmFailbackPolicyCreationInput{}

func (s InMageRcmFailbackPolicyCreationInput) MarshalJSON() ([]byte, error) {
	type wrapper InMageRcmFailbackPolicyCreationInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageRcmFailbackPolicyCreationInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageRcmFailbackPolicyCreationInput: %+v", err)
	}
	decoded["instanceType"] = "InMageRcmFailback"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageRcmFailbackPolicyCreationInput: %+v", err)
	}

	return encoded, nil
}
