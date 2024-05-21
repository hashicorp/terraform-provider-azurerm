package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EnableProtectionProviderSpecificInput = InMageEnableProtectionInput{}

type InMageEnableProtectionInput struct {
	DatastoreName      *string                   `json:"datastoreName,omitempty"`
	DiskExclusionInput *InMageDiskExclusionInput `json:"diskExclusionInput,omitempty"`
	DisksToInclude     *[]string                 `json:"disksToInclude,omitempty"`
	MasterTargetId     string                    `json:"masterTargetId"`
	MultiVMGroupId     string                    `json:"multiVmGroupId"`
	MultiVMGroupName   string                    `json:"multiVmGroupName"`
	ProcessServerId    string                    `json:"processServerId"`
	RetentionDrive     string                    `json:"retentionDrive"`
	RunAsAccountId     *string                   `json:"runAsAccountId,omitempty"`
	VMFriendlyName     *string                   `json:"vmFriendlyName,omitempty"`

	// Fields inherited from EnableProtectionProviderSpecificInput
}

var _ json.Marshaler = InMageEnableProtectionInput{}

func (s InMageEnableProtectionInput) MarshalJSON() ([]byte, error) {
	type wrapper InMageEnableProtectionInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageEnableProtectionInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageEnableProtectionInput: %+v", err)
	}
	decoded["instanceType"] = "InMage"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageEnableProtectionInput: %+v", err)
	}

	return encoded, nil
}
