package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AddDisksProviderSpecificInput = InMageRcmAddDisksInput{}

type InMageRcmAddDisksInput struct {
	Disks []InMageRcmDiskInput `json:"disks"`

	// Fields inherited from AddDisksProviderSpecificInput

	InstanceType string `json:"instanceType"`
}

func (s InMageRcmAddDisksInput) AddDisksProviderSpecificInput() BaseAddDisksProviderSpecificInputImpl {
	return BaseAddDisksProviderSpecificInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = InMageRcmAddDisksInput{}

func (s InMageRcmAddDisksInput) MarshalJSON() ([]byte, error) {
	type wrapper InMageRcmAddDisksInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageRcmAddDisksInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageRcmAddDisksInput: %+v", err)
	}

	decoded["instanceType"] = "InMageRcm"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageRcmAddDisksInput: %+v", err)
	}

	return encoded, nil
}
