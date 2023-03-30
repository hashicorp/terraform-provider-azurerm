package replicationnetworkmappings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ NetworkMappingFabricSpecificSettings = VMmToAzureNetworkMappingSettings{}

type VMmToAzureNetworkMappingSettings struct {

	// Fields inherited from NetworkMappingFabricSpecificSettings
}

var _ json.Marshaler = VMmToAzureNetworkMappingSettings{}

func (s VMmToAzureNetworkMappingSettings) MarshalJSON() ([]byte, error) {
	type wrapper VMmToAzureNetworkMappingSettings
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling VMmToAzureNetworkMappingSettings: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling VMmToAzureNetworkMappingSettings: %+v", err)
	}
	decoded["instanceType"] = "VmmToAzure"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling VMmToAzureNetworkMappingSettings: %+v", err)
	}

	return encoded, nil
}
