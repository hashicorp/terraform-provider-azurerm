package replicationpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ PolicyProviderSpecificDetails = InMageAzureV2PolicyDetails{}

type InMageAzureV2PolicyDetails struct {
	AppConsistentFrequencyInMinutes   *int64  `json:"appConsistentFrequencyInMinutes,omitempty"`
	CrashConsistentFrequencyInMinutes *int64  `json:"crashConsistentFrequencyInMinutes,omitempty"`
	MultiVMSyncStatus                 *string `json:"multiVmSyncStatus,omitempty"`
	RecoveryPointHistory              *int64  `json:"recoveryPointHistory,omitempty"`
	RecoveryPointThresholdInMinutes   *int64  `json:"recoveryPointThresholdInMinutes,omitempty"`

	// Fields inherited from PolicyProviderSpecificDetails

	InstanceType string `json:"instanceType"`
}

func (s InMageAzureV2PolicyDetails) PolicyProviderSpecificDetails() BasePolicyProviderSpecificDetailsImpl {
	return BasePolicyProviderSpecificDetailsImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = InMageAzureV2PolicyDetails{}

func (s InMageAzureV2PolicyDetails) MarshalJSON() ([]byte, error) {
	type wrapper InMageAzureV2PolicyDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageAzureV2PolicyDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageAzureV2PolicyDetails: %+v", err)
	}

	decoded["instanceType"] = "InMageAzureV2"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageAzureV2PolicyDetails: %+v", err)
	}

	return encoded, nil
}
