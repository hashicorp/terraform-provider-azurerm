package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ReverseReplicationProviderSpecificInput = InMageAzureV2ReprotectInput{}

type InMageAzureV2ReprotectInput struct {
	DisksToInclude      *[]string `json:"disksToInclude,omitempty"`
	LogStorageAccountId *string   `json:"logStorageAccountId,omitempty"`
	MasterTargetId      *string   `json:"masterTargetId,omitempty"`
	PolicyId            *string   `json:"policyId,omitempty"`
	ProcessServerId     *string   `json:"processServerId,omitempty"`
	RunAsAccountId      *string   `json:"runAsAccountId,omitempty"`
	StorageAccountId    *string   `json:"storageAccountId,omitempty"`

	// Fields inherited from ReverseReplicationProviderSpecificInput
}

var _ json.Marshaler = InMageAzureV2ReprotectInput{}

func (s InMageAzureV2ReprotectInput) MarshalJSON() ([]byte, error) {
	type wrapper InMageAzureV2ReprotectInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageAzureV2ReprotectInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageAzureV2ReprotectInput: %+v", err)
	}
	decoded["instanceType"] = "InMageAzureV2"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageAzureV2ReprotectInput: %+v", err)
	}

	return encoded, nil
}
