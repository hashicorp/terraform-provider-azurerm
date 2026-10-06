package replicationnetworkmappings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FabricSpecificCreateNetworkMappingInput = AzureToAzureCreateNetworkMappingInput{}

type AzureToAzureCreateNetworkMappingInput struct {
	PrimaryNetworkId string `json:"primaryNetworkId"`

	// Fields inherited from FabricSpecificCreateNetworkMappingInput

	InstanceType string `json:"instanceType"`
}

func (s AzureToAzureCreateNetworkMappingInput) FabricSpecificCreateNetworkMappingInput() BaseFabricSpecificCreateNetworkMappingInputImpl {
	return BaseFabricSpecificCreateNetworkMappingInputImpl{
		InstanceType: s.InstanceType,
	}
}

var _ json.Marshaler = AzureToAzureCreateNetworkMappingInput{}

func (s AzureToAzureCreateNetworkMappingInput) MarshalJSON() ([]byte, error) {
	type wrapper AzureToAzureCreateNetworkMappingInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureToAzureCreateNetworkMappingInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureToAzureCreateNetworkMappingInput: %+v", err)
	}

	decoded["instanceType"] = "AzureToAzure"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureToAzureCreateNetworkMappingInput: %+v", err)
	}

	return encoded, nil
}
