package replicationpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ PolicyProviderSpecificDetails = InMagePolicyDetails{}

type InMagePolicyDetails struct {
	AppConsistentFrequencyInMinutes *int64  `json:"appConsistentFrequencyInMinutes,omitempty"`
	MultiVMSyncStatus               *string `json:"multiVmSyncStatus,omitempty"`
	RecoveryPointHistory            *int64  `json:"recoveryPointHistory,omitempty"`
	RecoveryPointThresholdInMinutes *int64  `json:"recoveryPointThresholdInMinutes,omitempty"`

	// Fields inherited from PolicyProviderSpecificDetails
}

var _ json.Marshaler = InMagePolicyDetails{}

func (s InMagePolicyDetails) MarshalJSON() ([]byte, error) {
	type wrapper InMagePolicyDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMagePolicyDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMagePolicyDetails: %+v", err)
	}
	decoded["instanceType"] = "InMage"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMagePolicyDetails: %+v", err)
	}

	return encoded, nil
}
