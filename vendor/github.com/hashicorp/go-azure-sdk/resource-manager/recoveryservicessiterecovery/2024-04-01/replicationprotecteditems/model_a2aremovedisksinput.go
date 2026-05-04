package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ RemoveDisksProviderSpecificInput = A2ARemoveDisksInput{}

type A2ARemoveDisksInput struct {
	VMDisksUris       *[]string `json:"vmDisksUris,omitempty"`
	VMManagedDisksIds *[]string `json:"vmManagedDisksIds,omitempty"`

	// Fields inherited from RemoveDisksProviderSpecificInput

	InstanceType string `json:"instanceType"`
}

func (s A2ARemoveDisksInput) RemoveDisksProviderSpecificInput() BaseRemoveDisksProviderSpecificInputImpl {
	return BaseRemoveDisksProviderSpecificInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = A2ARemoveDisksInput{}

func (s A2ARemoveDisksInput) MarshalJSON() ([]byte, error) {
	type wrapper A2ARemoveDisksInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling A2ARemoveDisksInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling A2ARemoveDisksInput: %+v", err)
	}

	decoded["instanceType"] = "A2A"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling A2ARemoveDisksInput: %+v", err)
	}

	return encoded, nil
}
