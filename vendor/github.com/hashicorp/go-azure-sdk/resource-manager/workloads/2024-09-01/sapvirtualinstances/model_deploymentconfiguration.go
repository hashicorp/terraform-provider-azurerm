package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SAPConfiguration = DeploymentConfiguration{}

type DeploymentConfiguration struct {
	AppLocation                 *string                     `json:"appLocation,omitempty"`
	InfrastructureConfiguration InfrastructureConfiguration `json:"infrastructureConfiguration"`
	SoftwareConfiguration       SoftwareConfiguration       `json:"softwareConfiguration"`

	// Fields inherited from SAPConfiguration

	ConfigurationType SAPConfigurationType `json:"configurationType"`
}

func (s DeploymentConfiguration) SAPConfiguration() BaseSAPConfigurationImpl {
	return BaseSAPConfigurationImpl{
		ConfigurationType: s.ConfigurationType,
	}
}

var _ json.Marshaler = DeploymentConfiguration{}

func (s DeploymentConfiguration) MarshalJSON() ([]byte, error) {
	type wrapper DeploymentConfiguration
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DeploymentConfiguration: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DeploymentConfiguration: %+v", err)
	}

	decoded["configurationType"] = "Deployment"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DeploymentConfiguration: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &DeploymentConfiguration{}

func (s *DeploymentConfiguration) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AppLocation       *string              `json:"appLocation,omitempty"`
		ConfigurationType SAPConfigurationType `json:"configurationType"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AppLocation = decoded.AppLocation
	s.ConfigurationType = decoded.ConfigurationType

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling DeploymentConfiguration into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["infrastructureConfiguration"]; ok {
		impl, err := UnmarshalInfrastructureConfigurationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'InfrastructureConfiguration' for 'DeploymentConfiguration': %+v", err)
		}
		s.InfrastructureConfiguration = impl
	}

	if v, ok := temp["softwareConfiguration"]; ok {
		impl, err := UnmarshalSoftwareConfigurationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'SoftwareConfiguration' for 'DeploymentConfiguration': %+v", err)
		}
		s.SoftwareConfiguration = impl
	}

	return nil
}
