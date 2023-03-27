package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SAPConfiguration = DeploymentWithOSConfiguration{}

type DeploymentWithOSConfiguration struct {
	AppLocation                 *string                     `json:"appLocation,omitempty"`
	InfrastructureConfiguration InfrastructureConfiguration `json:"infrastructureConfiguration"`
	OsSapConfiguration          *OsSapConfiguration         `json:"osSapConfiguration,omitempty"`
	SoftwareConfiguration       SoftwareConfiguration       `json:"softwareConfiguration"`

	// Fields inherited from SAPConfiguration
}

var _ json.Marshaler = DeploymentWithOSConfiguration{}

func (s DeploymentWithOSConfiguration) MarshalJSON() ([]byte, error) {
	type wrapper DeploymentWithOSConfiguration
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DeploymentWithOSConfiguration: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DeploymentWithOSConfiguration: %+v", err)
	}
	decoded["configurationType"] = "DeploymentWithOSConfig"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DeploymentWithOSConfiguration: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &DeploymentWithOSConfiguration{}

func (s *DeploymentWithOSConfiguration) UnmarshalJSON(bytes []byte) error {
	type alias DeploymentWithOSConfiguration
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into DeploymentWithOSConfiguration: %+v", err)
	}

	s.AppLocation = decoded.AppLocation
	s.OsSapConfiguration = decoded.OsSapConfiguration

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling DeploymentWithOSConfiguration into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["infrastructureConfiguration"]; ok {
		impl, err := unmarshalInfrastructureConfigurationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'InfrastructureConfiguration' for 'DeploymentWithOSConfiguration': %+v", err)
		}
		s.InfrastructureConfiguration = impl
	}

	if v, ok := temp["softwareConfiguration"]; ok {
		impl, err := unmarshalSoftwareConfigurationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'SoftwareConfiguration' for 'DeploymentWithOSConfiguration': %+v", err)
		}
		s.SoftwareConfiguration = impl
	}
	return nil
}
