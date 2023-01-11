package replicationnetworkmappings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ NetworkMappingFabricSpecificSettings = VMmToVMmNetworkMappingSettings{}

type VMmToVMmNetworkMappingSettings struct {

	// Fields inherited from NetworkMappingFabricSpecificSettings
}

var _ json.Marshaler = VMmToVMmNetworkMappingSettings{}

func (s VMmToVMmNetworkMappingSettings) MarshalJSON() ([]byte, error) {
	type wrapper VMmToVMmNetworkMappingSettings
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling VMmToVMmNetworkMappingSettings: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling VMmToVMmNetworkMappingSettings: %+v", err)
	}
	decoded["instanceType"] = "VmmToVmm"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling VMmToVMmNetworkMappingSettings: %+v", err)
	}

	return encoded, nil
}
