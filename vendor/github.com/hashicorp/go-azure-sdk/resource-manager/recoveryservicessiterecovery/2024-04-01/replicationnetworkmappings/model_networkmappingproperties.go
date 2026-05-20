package replicationnetworkmappings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkMappingProperties struct {
	FabricSpecificSettings      NetworkMappingFabricSpecificSettings `json:"fabricSpecificSettings"`
	PrimaryFabricFriendlyName   *string                              `json:"primaryFabricFriendlyName,omitempty"`
	PrimaryNetworkFriendlyName  *string                              `json:"primaryNetworkFriendlyName,omitempty"`
	PrimaryNetworkId            *string                              `json:"primaryNetworkId,omitempty"`
	RecoveryFabricArmId         *string                              `json:"recoveryFabricArmId,omitempty"`
	RecoveryFabricFriendlyName  *string                              `json:"recoveryFabricFriendlyName,omitempty"`
	RecoveryNetworkFriendlyName *string                              `json:"recoveryNetworkFriendlyName,omitempty"`
	RecoveryNetworkId           *string                              `json:"recoveryNetworkId,omitempty"`
	State                       *string                              `json:"state,omitempty"`
}

var _ json.Unmarshaler = &NetworkMappingProperties{}

func (s *NetworkMappingProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		PrimaryFabricFriendlyName   *string `json:"primaryFabricFriendlyName,omitempty"`
		PrimaryNetworkFriendlyName  *string `json:"primaryNetworkFriendlyName,omitempty"`
		PrimaryNetworkId            *string `json:"primaryNetworkId,omitempty"`
		RecoveryFabricArmId         *string `json:"recoveryFabricArmId,omitempty"`
		RecoveryFabricFriendlyName  *string `json:"recoveryFabricFriendlyName,omitempty"`
		RecoveryNetworkFriendlyName *string `json:"recoveryNetworkFriendlyName,omitempty"`
		RecoveryNetworkId           *string `json:"recoveryNetworkId,omitempty"`
		State                       *string `json:"state,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.PrimaryFabricFriendlyName = decoded.PrimaryFabricFriendlyName
	s.PrimaryNetworkFriendlyName = decoded.PrimaryNetworkFriendlyName
	s.PrimaryNetworkId = decoded.PrimaryNetworkId
	s.RecoveryFabricArmId = decoded.RecoveryFabricArmId
	s.RecoveryFabricFriendlyName = decoded.RecoveryFabricFriendlyName
	s.RecoveryNetworkFriendlyName = decoded.RecoveryNetworkFriendlyName
	s.RecoveryNetworkId = decoded.RecoveryNetworkId
	s.State = decoded.State

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling NetworkMappingProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["fabricSpecificSettings"]; ok {
		impl, err := UnmarshalNetworkMappingFabricSpecificSettingsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'FabricSpecificSettings' for 'NetworkMappingProperties': %+v", err)
		}
		s.FabricSpecificSettings = impl
	}

	return nil
}
