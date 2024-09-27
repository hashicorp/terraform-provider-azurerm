package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SAPVirtualInstanceProperties struct {
	Configuration                     SAPConfiguration                     `json:"configuration"`
	Environment                       SAPEnvironmentType                   `json:"environment"`
	Errors                            *SAPVirtualInstanceError             `json:"errors,omitempty"`
	Health                            *SAPHealthState                      `json:"health,omitempty"`
	ManagedResourceGroupConfiguration *ManagedRGConfiguration              `json:"managedResourceGroupConfiguration,omitempty"`
	ProvisioningState                 *SapVirtualInstanceProvisioningState `json:"provisioningState,omitempty"`
	SapProduct                        SAPProductType                       `json:"sapProduct"`
	State                             *SAPVirtualInstanceState             `json:"state,omitempty"`
	Status                            *SAPVirtualInstanceStatus            `json:"status,omitempty"`
}

var _ json.Unmarshaler = &SAPVirtualInstanceProperties{}

func (s *SAPVirtualInstanceProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Environment                       SAPEnvironmentType                   `json:"environment"`
		Errors                            *SAPVirtualInstanceError             `json:"errors,omitempty"`
		Health                            *SAPHealthState                      `json:"health,omitempty"`
		ManagedResourceGroupConfiguration *ManagedRGConfiguration              `json:"managedResourceGroupConfiguration,omitempty"`
		ProvisioningState                 *SapVirtualInstanceProvisioningState `json:"provisioningState,omitempty"`
		SapProduct                        SAPProductType                       `json:"sapProduct"`
		State                             *SAPVirtualInstanceState             `json:"state,omitempty"`
		Status                            *SAPVirtualInstanceStatus            `json:"status,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Environment = decoded.Environment
	s.Errors = decoded.Errors
	s.Health = decoded.Health
	s.ManagedResourceGroupConfiguration = decoded.ManagedResourceGroupConfiguration
	s.ProvisioningState = decoded.ProvisioningState
	s.SapProduct = decoded.SapProduct
	s.State = decoded.State
	s.Status = decoded.Status

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SAPVirtualInstanceProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["configuration"]; ok {
		impl, err := UnmarshalSAPConfigurationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Configuration' for 'SAPVirtualInstanceProperties': %+v", err)
		}
		s.Configuration = impl
	}

	return nil
}
