package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ApplyRecoveryPointProviderSpecificInput = InMageRcmApplyRecoveryPointInput{}

type InMageRcmApplyRecoveryPointInput struct {
	RecoveryPointId string `json:"recoveryPointId"`

	// Fields inherited from ApplyRecoveryPointProviderSpecificInput

	InstanceType string `json:"instanceType"`
}

func (s InMageRcmApplyRecoveryPointInput) ApplyRecoveryPointProviderSpecificInput() BaseApplyRecoveryPointProviderSpecificInputImpl {
	return BaseApplyRecoveryPointProviderSpecificInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = InMageRcmApplyRecoveryPointInput{}

func (s InMageRcmApplyRecoveryPointInput) MarshalJSON() ([]byte, error) {
	type wrapper InMageRcmApplyRecoveryPointInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageRcmApplyRecoveryPointInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageRcmApplyRecoveryPointInput: %+v", err)
	}

	decoded["instanceType"] = "InMageRcm"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageRcmApplyRecoveryPointInput: %+v", err)
	}

	return encoded, nil
}
