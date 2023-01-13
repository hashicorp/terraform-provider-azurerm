package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ReverseReplicationProviderSpecificInput = InMageReprotectInput{}

type InMageReprotectInput struct {
	DatastoreName      *string                   `json:"datastoreName,omitempty"`
	DiskExclusionInput *InMageDiskExclusionInput `json:"diskExclusionInput,omitempty"`
	DisksToInclude     *[]string                 `json:"disksToInclude,omitempty"`
	MasterTargetId     string                    `json:"masterTargetId"`
	ProcessServerId    string                    `json:"processServerId"`
	ProfileId          string                    `json:"profileId"`
	RetentionDrive     string                    `json:"retentionDrive"`
	RunAsAccountId     *string                   `json:"runAsAccountId,omitempty"`

	// Fields inherited from ReverseReplicationProviderSpecificInput
}

var _ json.Marshaler = InMageReprotectInput{}

func (s InMageReprotectInput) MarshalJSON() ([]byte, error) {
	type wrapper InMageReprotectInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageReprotectInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageReprotectInput: %+v", err)
	}
	decoded["instanceType"] = "InMage"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageReprotectInput: %+v", err)
	}

	return encoded, nil
}
