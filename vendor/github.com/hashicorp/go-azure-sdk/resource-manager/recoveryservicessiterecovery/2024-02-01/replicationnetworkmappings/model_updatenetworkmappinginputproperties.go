package replicationnetworkmappings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateNetworkMappingInputProperties struct {
	FabricSpecificDetails FabricSpecificUpdateNetworkMappingInput `json:"fabricSpecificDetails"`
	RecoveryFabricName    *string                                 `json:"recoveryFabricName,omitempty"`
	RecoveryNetworkId     *string                                 `json:"recoveryNetworkId,omitempty"`
}

var _ json.Unmarshaler = &UpdateNetworkMappingInputProperties{}

func (s *UpdateNetworkMappingInputProperties) UnmarshalJSON(bytes []byte) error {
	type alias UpdateNetworkMappingInputProperties
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into UpdateNetworkMappingInputProperties: %+v", err)
	}

	s.RecoveryFabricName = decoded.RecoveryFabricName
	s.RecoveryNetworkId = decoded.RecoveryNetworkId

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling UpdateNetworkMappingInputProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["fabricSpecificDetails"]; ok {
		impl, err := unmarshalFabricSpecificUpdateNetworkMappingInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'FabricSpecificDetails' for 'UpdateNetworkMappingInputProperties': %+v", err)
		}
		s.FabricSpecificDetails = impl
	}
	return nil
}
