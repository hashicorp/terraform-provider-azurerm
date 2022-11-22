package replicationnetworkmappings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FabricSpecificUpdateNetworkMappingInput = VmmToVmmUpdateNetworkMappingInput{}

type VmmToVmmUpdateNetworkMappingInput struct {

	// Fields inherited from FabricSpecificUpdateNetworkMappingInput
}

var _ json.Marshaler = VmmToVmmUpdateNetworkMappingInput{}

func (s VmmToVmmUpdateNetworkMappingInput) MarshalJSON() ([]byte, error) {
	type wrapper VmmToVmmUpdateNetworkMappingInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling VmmToVmmUpdateNetworkMappingInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling VmmToVmmUpdateNetworkMappingInput: %+v", err)
	}
	decoded["instanceType"] = "VmmToVmm"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling VmmToVmmUpdateNetworkMappingInput: %+v", err)
	}

	return encoded, nil
}
