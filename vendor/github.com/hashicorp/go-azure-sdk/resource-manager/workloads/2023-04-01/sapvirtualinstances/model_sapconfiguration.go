package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SAPConfiguration interface {
	SAPConfiguration() BaseSAPConfigurationImpl
}

var _ SAPConfiguration = BaseSAPConfigurationImpl{}

type BaseSAPConfigurationImpl struct {
	ConfigurationType SAPConfigurationType `json:"configurationType"`
}

func (s BaseSAPConfigurationImpl) SAPConfiguration() BaseSAPConfigurationImpl {
	return s
}

var _ SAPConfiguration = RawSAPConfigurationImpl{}

// RawSAPConfigurationImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawSAPConfigurationImpl struct {
	sAPConfiguration BaseSAPConfigurationImpl
	Type             string
	Values           map[string]interface{}
}

func (s RawSAPConfigurationImpl) SAPConfiguration() BaseSAPConfigurationImpl {
	return s.sAPConfiguration
}

func UnmarshalSAPConfigurationImplementation(input []byte) (SAPConfiguration, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling SAPConfiguration into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["configurationType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Deployment") {
		var out DeploymentConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DeploymentConfiguration: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DeploymentWithOSConfig") {
		var out DeploymentWithOSConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DeploymentWithOSConfiguration: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Discovery") {
		var out DiscoveryConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DiscoveryConfiguration: %+v", err)
		}
		return out, nil
	}

	var parent BaseSAPConfigurationImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseSAPConfigurationImpl: %+v", err)
	}

	return RawSAPConfigurationImpl{
		sAPConfiguration: parent,
		Type:             value,
		Values:           temp,
	}, nil

}
