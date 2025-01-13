package replicationpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ PolicyProviderSpecificDetails = InMageBasePolicyDetails{}

type InMageBasePolicyDetails struct {
	AppConsistentFrequencyInMinutes *int64  `json:"appConsistentFrequencyInMinutes,omitempty"`
	MultiVMSyncStatus               *string `json:"multiVmSyncStatus,omitempty"`
	RecoveryPointHistory            *int64  `json:"recoveryPointHistory,omitempty"`
	RecoveryPointThresholdInMinutes *int64  `json:"recoveryPointThresholdInMinutes,omitempty"`

	// Fields inherited from PolicyProviderSpecificDetails

	InstanceType string `json:"instanceType"`
}

func (s InMageBasePolicyDetails) PolicyProviderSpecificDetails() BasePolicyProviderSpecificDetailsImpl {
	return BasePolicyProviderSpecificDetailsImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = InMageBasePolicyDetails{}

func (s InMageBasePolicyDetails) MarshalJSON() ([]byte, error) {
	type wrapper InMageBasePolicyDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageBasePolicyDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageBasePolicyDetails: %+v", err)
	}

	decoded["instanceType"] = "InMageBasePolicyDetails"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageBasePolicyDetails: %+v", err)
	}

	return encoded, nil
}
