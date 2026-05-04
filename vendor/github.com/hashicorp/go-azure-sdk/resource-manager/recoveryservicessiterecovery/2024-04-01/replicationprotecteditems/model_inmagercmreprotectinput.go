package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ReverseReplicationProviderSpecificInput = InMageRcmReprotectInput{}

type InMageRcmReprotectInput struct {
	DatastoreName       string  `json:"datastoreName"`
	LogStorageAccountId string  `json:"logStorageAccountId"`
	PolicyId            *string `json:"policyId,omitempty"`
	ReprotectAgentId    string  `json:"reprotectAgentId"`

	// Fields inherited from ReverseReplicationProviderSpecificInput

	InstanceType string `json:"instanceType"`
}

func (s InMageRcmReprotectInput) ReverseReplicationProviderSpecificInput() BaseReverseReplicationProviderSpecificInputImpl {
	return BaseReverseReplicationProviderSpecificInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = InMageRcmReprotectInput{}

func (s InMageRcmReprotectInput) MarshalJSON() ([]byte, error) {
	type wrapper InMageRcmReprotectInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageRcmReprotectInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageRcmReprotectInput: %+v", err)
	}

	decoded["instanceType"] = "InMageRcm"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageRcmReprotectInput: %+v", err)
	}

	return encoded, nil
}
