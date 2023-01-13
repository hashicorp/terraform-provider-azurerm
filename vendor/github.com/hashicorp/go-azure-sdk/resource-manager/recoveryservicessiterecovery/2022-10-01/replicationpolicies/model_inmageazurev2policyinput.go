package replicationpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ PolicyProviderSpecificInput = InMageAzureV2PolicyInput{}

type InMageAzureV2PolicyInput struct {
	AppConsistentFrequencyInMinutes   *int64               `json:"appConsistentFrequencyInMinutes,omitempty"`
	CrashConsistentFrequencyInMinutes *int64               `json:"crashConsistentFrequencyInMinutes,omitempty"`
	MultiVMSyncStatus                 SetMultiVMSyncStatus `json:"multiVmSyncStatus"`
	RecoveryPointHistory              *int64               `json:"recoveryPointHistory,omitempty"`
	RecoveryPointThresholdInMinutes   *int64               `json:"recoveryPointThresholdInMinutes,omitempty"`

	// Fields inherited from PolicyProviderSpecificInput
}

var _ json.Marshaler = InMageAzureV2PolicyInput{}

func (s InMageAzureV2PolicyInput) MarshalJSON() ([]byte, error) {
	type wrapper InMageAzureV2PolicyInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageAzureV2PolicyInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageAzureV2PolicyInput: %+v", err)
	}
	decoded["instanceType"] = "InMageAzureV2"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageAzureV2PolicyInput: %+v", err)
	}

	return encoded, nil
}
