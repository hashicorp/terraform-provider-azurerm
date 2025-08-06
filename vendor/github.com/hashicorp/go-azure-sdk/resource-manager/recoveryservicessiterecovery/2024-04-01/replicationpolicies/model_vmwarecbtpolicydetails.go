package replicationpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ PolicyProviderSpecificDetails = VMwareCbtPolicyDetails{}

type VMwareCbtPolicyDetails struct {
	AppConsistentFrequencyInMinutes   *int64 `json:"appConsistentFrequencyInMinutes,omitempty"`
	CrashConsistentFrequencyInMinutes *int64 `json:"crashConsistentFrequencyInMinutes,omitempty"`
	RecoveryPointHistoryInMinutes     *int64 `json:"recoveryPointHistoryInMinutes,omitempty"`

	// Fields inherited from PolicyProviderSpecificDetails

	InstanceType string `json:"instanceType"`
}

func (s VMwareCbtPolicyDetails) PolicyProviderSpecificDetails() BasePolicyProviderSpecificDetailsImpl {
	return BasePolicyProviderSpecificDetailsImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = VMwareCbtPolicyDetails{}

func (s VMwareCbtPolicyDetails) MarshalJSON() ([]byte, error) {
	type wrapper VMwareCbtPolicyDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling VMwareCbtPolicyDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling VMwareCbtPolicyDetails: %+v", err)
	}

	decoded["instanceType"] = "VMwareCbt"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling VMwareCbtPolicyDetails: %+v", err)
	}

	return encoded, nil
}
