package replicationpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ PolicyProviderSpecificDetails = InMageRcmPolicyDetails{}

type InMageRcmPolicyDetails struct {
	AppConsistentFrequencyInMinutes   *int64  `json:"appConsistentFrequencyInMinutes,omitempty"`
	CrashConsistentFrequencyInMinutes *int64  `json:"crashConsistentFrequencyInMinutes,omitempty"`
	EnableMultiVMSync                 *string `json:"enableMultiVmSync,omitempty"`
	RecoveryPointHistoryInMinutes     *int64  `json:"recoveryPointHistoryInMinutes,omitempty"`

	// Fields inherited from PolicyProviderSpecificDetails
}

var _ json.Marshaler = InMageRcmPolicyDetails{}

func (s InMageRcmPolicyDetails) MarshalJSON() ([]byte, error) {
	type wrapper InMageRcmPolicyDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageRcmPolicyDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageRcmPolicyDetails: %+v", err)
	}
	decoded["instanceType"] = "InMageRcm"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageRcmPolicyDetails: %+v", err)
	}

	return encoded, nil
}
