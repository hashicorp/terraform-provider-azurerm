package pools

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PoolUpdateProperties struct {
	AgentProfile               AgentProfile        `json:"agentProfile"`
	DevCenterProjectResourceId *string             `json:"devCenterProjectResourceId,omitempty"`
	FabricProfile              FabricProfile       `json:"fabricProfile"`
	MaximumConcurrency         *int64              `json:"maximumConcurrency,omitempty"`
	OrganizationProfile        OrganizationProfile `json:"organizationProfile"`
	ProvisioningState          *ProvisioningState  `json:"provisioningState,omitempty"`
}

var _ json.Unmarshaler = &PoolUpdateProperties{}

func (s *PoolUpdateProperties) UnmarshalJSON(bytes []byte) error {
	type alias PoolUpdateProperties
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into PoolUpdateProperties: %+v", err)
	}

	s.DevCenterProjectResourceId = decoded.DevCenterProjectResourceId
	s.MaximumConcurrency = decoded.MaximumConcurrency
	s.ProvisioningState = decoded.ProvisioningState

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling PoolUpdateProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["agentProfile"]; ok {
		impl, err := UnmarshalAgentProfileImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'AgentProfile' for 'PoolUpdateProperties': %+v", err)
		}
		s.AgentProfile = impl
	}

	if v, ok := temp["fabricProfile"]; ok {
		impl, err := UnmarshalFabricProfileImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'FabricProfile' for 'PoolUpdateProperties': %+v", err)
		}
		s.FabricProfile = impl
	}

	if v, ok := temp["organizationProfile"]; ok {
		impl, err := UnmarshalOrganizationProfileImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'OrganizationProfile' for 'PoolUpdateProperties': %+v", err)
		}
		s.OrganizationProfile = impl
	}
	return nil
}
