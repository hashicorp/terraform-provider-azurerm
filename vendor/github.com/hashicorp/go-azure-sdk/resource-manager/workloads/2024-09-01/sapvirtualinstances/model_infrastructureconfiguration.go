package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InfrastructureConfiguration interface {
	InfrastructureConfiguration() BaseInfrastructureConfigurationImpl
}

var _ InfrastructureConfiguration = BaseInfrastructureConfigurationImpl{}

type BaseInfrastructureConfigurationImpl struct {
	AppResourceGroup string            `json:"appResourceGroup"`
	DeploymentType   SAPDeploymentType `json:"deploymentType"`
}

func (s BaseInfrastructureConfigurationImpl) InfrastructureConfiguration() BaseInfrastructureConfigurationImpl {
	return s
}

var _ InfrastructureConfiguration = RawInfrastructureConfigurationImpl{}

// RawInfrastructureConfigurationImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawInfrastructureConfigurationImpl struct {
	infrastructureConfiguration BaseInfrastructureConfigurationImpl
	Type                        string
	Values                      map[string]interface{}
}

func (s RawInfrastructureConfigurationImpl) InfrastructureConfiguration() BaseInfrastructureConfigurationImpl {
	return s.infrastructureConfiguration
}

func (s RawInfrastructureConfigurationImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalInfrastructureConfigurationImplementation(input []byte) (InfrastructureConfiguration, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling InfrastructureConfiguration into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["deploymentType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "SingleServer") {
		var out SingleServerConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SingleServerConfiguration: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ThreeTier") {
		var out ThreeTierConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ThreeTierConfiguration: %+v", err)
		}
		return out, nil
	}

	var parent BaseInfrastructureConfigurationImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseInfrastructureConfigurationImpl: %+v", err)
	}

	return RawInfrastructureConfigurationImpl{
		infrastructureConfiguration: parent,
		Type:                        value,
		Values:                      temp,
	}, nil

}
