package replicationnetworkmappings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FabricSpecificUpdateNetworkMappingInput = AzureToAzureUpdateNetworkMappingInput{}

type AzureToAzureUpdateNetworkMappingInput struct {
	PrimaryNetworkId *string `json:"primaryNetworkId,omitempty"`

	// Fields inherited from FabricSpecificUpdateNetworkMappingInput
}

var _ json.Marshaler = AzureToAzureUpdateNetworkMappingInput{}

func (s AzureToAzureUpdateNetworkMappingInput) MarshalJSON() ([]byte, error) {
	type wrapper AzureToAzureUpdateNetworkMappingInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureToAzureUpdateNetworkMappingInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureToAzureUpdateNetworkMappingInput: %+v", err)
	}
	decoded["instanceType"] = "AzureToAzure"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureToAzureUpdateNetworkMappingInput: %+v", err)
	}

	return encoded, nil
}
