package replicationpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ PolicyProviderSpecificDetails = InMageRcmFailbackPolicyDetails{}

type InMageRcmFailbackPolicyDetails struct {
	AppConsistentFrequencyInMinutes   *int64 `json:"appConsistentFrequencyInMinutes,omitempty"`
	CrashConsistentFrequencyInMinutes *int64 `json:"crashConsistentFrequencyInMinutes,omitempty"`

	// Fields inherited from PolicyProviderSpecificDetails

	InstanceType string `json:"instanceType"`
}

func (s InMageRcmFailbackPolicyDetails) PolicyProviderSpecificDetails() BasePolicyProviderSpecificDetailsImpl {
	return BasePolicyProviderSpecificDetailsImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = InMageRcmFailbackPolicyDetails{}

func (s InMageRcmFailbackPolicyDetails) MarshalJSON() ([]byte, error) {
	type wrapper InMageRcmFailbackPolicyDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageRcmFailbackPolicyDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageRcmFailbackPolicyDetails: %+v", err)
	}

	decoded["instanceType"] = "InMageRcmFailback"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageRcmFailbackPolicyDetails: %+v", err)
	}

	return encoded, nil
}
