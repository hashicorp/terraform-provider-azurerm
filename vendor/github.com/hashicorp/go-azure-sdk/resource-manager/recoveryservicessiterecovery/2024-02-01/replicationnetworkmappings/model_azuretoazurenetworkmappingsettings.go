package replicationnetworkmappings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ NetworkMappingFabricSpecificSettings = AzureToAzureNetworkMappingSettings{}

type AzureToAzureNetworkMappingSettings struct {
	PrimaryFabricLocation  *string `json:"primaryFabricLocation,omitempty"`
	RecoveryFabricLocation *string `json:"recoveryFabricLocation,omitempty"`

	// Fields inherited from NetworkMappingFabricSpecificSettings
}

var _ json.Marshaler = AzureToAzureNetworkMappingSettings{}

func (s AzureToAzureNetworkMappingSettings) MarshalJSON() ([]byte, error) {
	type wrapper AzureToAzureNetworkMappingSettings
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureToAzureNetworkMappingSettings: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureToAzureNetworkMappingSettings: %+v", err)
	}
	decoded["instanceType"] = "AzureToAzure"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureToAzureNetworkMappingSettings: %+v", err)
	}

	return encoded, nil
}
