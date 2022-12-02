package replicationnetworkmappings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FabricSpecificUpdateNetworkMappingInput = VmmToAzureUpdateNetworkMappingInput{}

type VmmToAzureUpdateNetworkMappingInput struct {

	// Fields inherited from FabricSpecificUpdateNetworkMappingInput
}

var _ json.Marshaler = VmmToAzureUpdateNetworkMappingInput{}

func (s VmmToAzureUpdateNetworkMappingInput) MarshalJSON() ([]byte, error) {
	type wrapper VmmToAzureUpdateNetworkMappingInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling VmmToAzureUpdateNetworkMappingInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling VmmToAzureUpdateNetworkMappingInput: %+v", err)
	}
	decoded["instanceType"] = "VmmToAzure"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling VmmToAzureUpdateNetworkMappingInput: %+v", err)
	}

	return encoded, nil
}
