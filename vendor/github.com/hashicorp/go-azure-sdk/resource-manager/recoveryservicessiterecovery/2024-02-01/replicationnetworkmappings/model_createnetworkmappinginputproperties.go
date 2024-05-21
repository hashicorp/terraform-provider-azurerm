package replicationnetworkmappings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CreateNetworkMappingInputProperties struct {
	FabricSpecificDetails FabricSpecificCreateNetworkMappingInput `json:"fabricSpecificDetails"`
	RecoveryFabricName    *string                                 `json:"recoveryFabricName,omitempty"`
	RecoveryNetworkId     string                                  `json:"recoveryNetworkId"`
}

var _ json.Unmarshaler = &CreateNetworkMappingInputProperties{}

func (s *CreateNetworkMappingInputProperties) UnmarshalJSON(bytes []byte) error {
	type alias CreateNetworkMappingInputProperties
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into CreateNetworkMappingInputProperties: %+v", err)
	}

	s.RecoveryFabricName = decoded.RecoveryFabricName
	s.RecoveryNetworkId = decoded.RecoveryNetworkId

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling CreateNetworkMappingInputProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["fabricSpecificDetails"]; ok {
		impl, err := unmarshalFabricSpecificCreateNetworkMappingInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'FabricSpecificDetails' for 'CreateNetworkMappingInputProperties': %+v", err)
		}
		s.FabricSpecificDetails = impl
	}
	return nil
}
