package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ApplyRecoveryPointProviderSpecificInput = InMageAzureV2ApplyRecoveryPointInput{}

type InMageAzureV2ApplyRecoveryPointInput struct {

	// Fields inherited from ApplyRecoveryPointProviderSpecificInput

	InstanceType string `json:"instanceType"`
}

func (s InMageAzureV2ApplyRecoveryPointInput) ApplyRecoveryPointProviderSpecificInput() BaseApplyRecoveryPointProviderSpecificInputImpl {
	return BaseApplyRecoveryPointProviderSpecificInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = InMageAzureV2ApplyRecoveryPointInput{}

func (s InMageAzureV2ApplyRecoveryPointInput) MarshalJSON() ([]byte, error) {
	type wrapper InMageAzureV2ApplyRecoveryPointInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageAzureV2ApplyRecoveryPointInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageAzureV2ApplyRecoveryPointInput: %+v", err)
	}

	decoded["instanceType"] = "InMageAzureV2"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageAzureV2ApplyRecoveryPointInput: %+v", err)
	}

	return encoded, nil
}
