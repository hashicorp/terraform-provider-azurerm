package replicationprotectioncontainermappings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateProtectionContainerMappingInputProperties struct {
	ProviderSpecificInput ReplicationProviderSpecificUpdateContainerMappingInput `json:"providerSpecificInput"`
}

var _ json.Unmarshaler = &UpdateProtectionContainerMappingInputProperties{}

func (s *UpdateProtectionContainerMappingInputProperties) UnmarshalJSON(bytes []byte) error {

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling UpdateProtectionContainerMappingInputProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["providerSpecificInput"]; ok {
		impl, err := UnmarshalReplicationProviderSpecificUpdateContainerMappingInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ProviderSpecificInput' for 'UpdateProtectionContainerMappingInputProperties': %+v", err)
		}
		s.ProviderSpecificInput = impl
	}

	return nil
}
