package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AddDisksProviderSpecificInput = A2AAddDisksInput{}

type A2AAddDisksInput struct {
	VMDisks        *[]A2AVMDiskInputDetails        `json:"vmDisks,omitempty"`
	VMManagedDisks *[]A2AVMManagedDiskInputDetails `json:"vmManagedDisks,omitempty"`

	// Fields inherited from AddDisksProviderSpecificInput

	InstanceType string `json:"instanceType"`
}

func (s A2AAddDisksInput) AddDisksProviderSpecificInput() BaseAddDisksProviderSpecificInputImpl {
	return BaseAddDisksProviderSpecificInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = A2AAddDisksInput{}

func (s A2AAddDisksInput) MarshalJSON() ([]byte, error) {
	type wrapper A2AAddDisksInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling A2AAddDisksInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling A2AAddDisksInput: %+v", err)
	}

	decoded["instanceType"] = "A2A"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling A2AAddDisksInput: %+v", err)
	}

	return encoded, nil
}
