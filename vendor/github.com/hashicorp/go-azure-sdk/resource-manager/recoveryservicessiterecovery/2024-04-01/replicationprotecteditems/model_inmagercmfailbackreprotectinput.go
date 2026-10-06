package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ReverseReplicationProviderSpecificInput = InMageRcmFailbackReprotectInput{}

type InMageRcmFailbackReprotectInput struct {
	PolicyId        string  `json:"policyId"`
	ProcessServerId string  `json:"processServerId"`
	RunAsAccountId  *string `json:"runAsAccountId,omitempty"`

	// Fields inherited from ReverseReplicationProviderSpecificInput

	InstanceType string `json:"instanceType"`
}

func (s InMageRcmFailbackReprotectInput) ReverseReplicationProviderSpecificInput() BaseReverseReplicationProviderSpecificInputImpl {
	return BaseReverseReplicationProviderSpecificInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = InMageRcmFailbackReprotectInput{}

func (s InMageRcmFailbackReprotectInput) MarshalJSON() ([]byte, error) {
	type wrapper InMageRcmFailbackReprotectInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageRcmFailbackReprotectInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageRcmFailbackReprotectInput: %+v", err)
	}

	decoded["instanceType"] = "InMageRcmFailback"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageRcmFailbackReprotectInput: %+v", err)
	}

	return encoded, nil
}
