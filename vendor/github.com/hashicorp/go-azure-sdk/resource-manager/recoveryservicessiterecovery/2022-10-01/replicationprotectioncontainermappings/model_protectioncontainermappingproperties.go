package replicationprotectioncontainermappings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProtectionContainerMappingProperties struct {
	Health                                *string                                           `json:"health,omitempty"`
	HealthErrorDetails                    *[]HealthError                                    `json:"healthErrorDetails,omitempty"`
	PolicyFriendlyName                    *string                                           `json:"policyFriendlyName,omitempty"`
	PolicyId                              *string                                           `json:"policyId,omitempty"`
	ProviderSpecificDetails               ProtectionContainerMappingProviderSpecificDetails `json:"providerSpecificDetails"`
	SourceFabricFriendlyName              *string                                           `json:"sourceFabricFriendlyName,omitempty"`
	SourceProtectionContainerFriendlyName *string                                           `json:"sourceProtectionContainerFriendlyName,omitempty"`
	State                                 *string                                           `json:"state,omitempty"`
	TargetFabricFriendlyName              *string                                           `json:"targetFabricFriendlyName,omitempty"`
	TargetProtectionContainerFriendlyName *string                                           `json:"targetProtectionContainerFriendlyName,omitempty"`
	TargetProtectionContainerId           *string                                           `json:"targetProtectionContainerId,omitempty"`
}

var _ json.Unmarshaler = &ProtectionContainerMappingProperties{}

func (s *ProtectionContainerMappingProperties) UnmarshalJSON(bytes []byte) error {
	type alias ProtectionContainerMappingProperties
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into ProtectionContainerMappingProperties: %+v", err)
	}

	s.Health = decoded.Health
	s.HealthErrorDetails = decoded.HealthErrorDetails
	s.PolicyFriendlyName = decoded.PolicyFriendlyName
	s.PolicyId = decoded.PolicyId
	s.SourceFabricFriendlyName = decoded.SourceFabricFriendlyName
	s.SourceProtectionContainerFriendlyName = decoded.SourceProtectionContainerFriendlyName
	s.State = decoded.State
	s.TargetFabricFriendlyName = decoded.TargetFabricFriendlyName
	s.TargetProtectionContainerFriendlyName = decoded.TargetProtectionContainerFriendlyName
	s.TargetProtectionContainerId = decoded.TargetProtectionContainerId

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ProtectionContainerMappingProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["providerSpecificDetails"]; ok {
		impl, err := unmarshalProtectionContainerMappingProviderSpecificDetailsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ProviderSpecificDetails' for 'ProtectionContainerMappingProperties': %+v", err)
		}
		s.ProviderSpecificDetails = impl
	}
	return nil
}
