package replicationnetworkmappings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FabricSpecificCreateNetworkMappingInput = VMmToAzureCreateNetworkMappingInput{}

type VMmToAzureCreateNetworkMappingInput struct {

	// Fields inherited from FabricSpecificCreateNetworkMappingInput
}

var _ json.Marshaler = VMmToAzureCreateNetworkMappingInput{}

func (s VMmToAzureCreateNetworkMappingInput) MarshalJSON() ([]byte, error) {
	type wrapper VMmToAzureCreateNetworkMappingInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling VMmToAzureCreateNetworkMappingInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling VMmToAzureCreateNetworkMappingInput: %+v", err)
	}
	decoded["instanceType"] = "VmmToAzure"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling VMmToAzureCreateNetworkMappingInput: %+v", err)
	}

	return encoded, nil
}
