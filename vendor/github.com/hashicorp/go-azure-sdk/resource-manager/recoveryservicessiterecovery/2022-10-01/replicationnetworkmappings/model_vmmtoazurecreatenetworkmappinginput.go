package replicationnetworkmappings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FabricSpecificCreateNetworkMappingInput = VmmToAzureCreateNetworkMappingInput{}

type VmmToAzureCreateNetworkMappingInput struct {

	// Fields inherited from FabricSpecificCreateNetworkMappingInput
}

var _ json.Marshaler = VmmToAzureCreateNetworkMappingInput{}

func (s VmmToAzureCreateNetworkMappingInput) MarshalJSON() ([]byte, error) {
	type wrapper VmmToAzureCreateNetworkMappingInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling VmmToAzureCreateNetworkMappingInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling VmmToAzureCreateNetworkMappingInput: %+v", err)
	}
	decoded["instanceType"] = "VmmToAzure"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling VmmToAzureCreateNetworkMappingInput: %+v", err)
	}

	return encoded, nil
}
