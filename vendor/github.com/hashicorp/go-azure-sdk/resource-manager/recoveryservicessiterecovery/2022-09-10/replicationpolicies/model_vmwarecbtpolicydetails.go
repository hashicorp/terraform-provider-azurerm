package replicationpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ PolicyProviderSpecificDetails = VmwareCbtPolicyDetails{}

type VmwareCbtPolicyDetails struct {
	AppConsistentFrequencyInMinutes   *int64 `json:"appConsistentFrequencyInMinutes,omitempty"`
	CrashConsistentFrequencyInMinutes *int64 `json:"crashConsistentFrequencyInMinutes,omitempty"`
	RecoveryPointHistoryInMinutes     *int64 `json:"recoveryPointHistoryInMinutes,omitempty"`

	// Fields inherited from PolicyProviderSpecificDetails
}

var _ json.Marshaler = VmwareCbtPolicyDetails{}

func (s VmwareCbtPolicyDetails) MarshalJSON() ([]byte, error) {
	type wrapper VmwareCbtPolicyDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling VmwareCbtPolicyDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling VmwareCbtPolicyDetails: %+v", err)
	}
	decoded["instanceType"] = "VMwareCbt"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling VmwareCbtPolicyDetails: %+v", err)
	}

	return encoded, nil
}
