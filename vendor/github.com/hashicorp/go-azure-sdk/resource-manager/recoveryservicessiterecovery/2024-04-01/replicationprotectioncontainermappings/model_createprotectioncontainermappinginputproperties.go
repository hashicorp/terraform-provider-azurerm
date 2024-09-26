package replicationprotectioncontainermappings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CreateProtectionContainerMappingInputProperties struct {
	PolicyId                    *string                                          `json:"policyId,omitempty"`
	ProviderSpecificInput       ReplicationProviderSpecificContainerMappingInput `json:"providerSpecificInput"`
	TargetProtectionContainerId *string                                          `json:"targetProtectionContainerId,omitempty"`
}

var _ json.Unmarshaler = &CreateProtectionContainerMappingInputProperties{}

func (s *CreateProtectionContainerMappingInputProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		PolicyId                    *string `json:"policyId,omitempty"`
		TargetProtectionContainerId *string `json:"targetProtectionContainerId,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.PolicyId = decoded.PolicyId
	s.TargetProtectionContainerId = decoded.TargetProtectionContainerId

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling CreateProtectionContainerMappingInputProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["providerSpecificInput"]; ok {
		impl, err := UnmarshalReplicationProviderSpecificContainerMappingInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ProviderSpecificInput' for 'CreateProtectionContainerMappingInputProperties': %+v", err)
		}
		s.ProviderSpecificInput = impl
	}

	return nil
}
