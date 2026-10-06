package replicationpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ PolicyProviderSpecificInput = InMagePolicyInput{}

type InMagePolicyInput struct {
	AppConsistentFrequencyInMinutes *int64               `json:"appConsistentFrequencyInMinutes,omitempty"`
	MultiVMSyncStatus               SetMultiVMSyncStatus `json:"multiVmSyncStatus"`
	RecoveryPointHistory            *int64               `json:"recoveryPointHistory,omitempty"`
	RecoveryPointThresholdInMinutes *int64               `json:"recoveryPointThresholdInMinutes,omitempty"`

	// Fields inherited from PolicyProviderSpecificInput

	InstanceType string `json:"instanceType"`
}

func (s InMagePolicyInput) PolicyProviderSpecificInput() BasePolicyProviderSpecificInputImpl {
	return BasePolicyProviderSpecificInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = InMagePolicyInput{}

func (s InMagePolicyInput) MarshalJSON() ([]byte, error) {
	type wrapper InMagePolicyInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMagePolicyInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMagePolicyInput: %+v", err)
	}

	decoded["instanceType"] = "InMage"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMagePolicyInput: %+v", err)
	}

	return encoded, nil
}
