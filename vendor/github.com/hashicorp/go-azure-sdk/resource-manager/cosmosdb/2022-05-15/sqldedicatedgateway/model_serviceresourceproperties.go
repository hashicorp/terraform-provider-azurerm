package sqldedicatedgateway

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceResourceProperties interface {
	ServiceResourceProperties() BaseServiceResourcePropertiesImpl
}

var _ ServiceResourceProperties = BaseServiceResourcePropertiesImpl{}

type BaseServiceResourcePropertiesImpl struct {
	CreationTime  *string        `json:"creationTime,omitempty"`
	InstanceCount *int64         `json:"instanceCount,omitempty"`
	InstanceSize  *ServiceSize   `json:"instanceSize,omitempty"`
	ServiceType   ServiceType    `json:"serviceType"`
	Status        *ServiceStatus `json:"status,omitempty"`
}

func (s BaseServiceResourcePropertiesImpl) ServiceResourceProperties() BaseServiceResourcePropertiesImpl {
	return s
}

var _ ServiceResourceProperties = RawServiceResourcePropertiesImpl{}

// RawServiceResourcePropertiesImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawServiceResourcePropertiesImpl struct {
	serviceResourceProperties BaseServiceResourcePropertiesImpl
	Type                      string
	Values                    map[string]interface{}
}

func (s RawServiceResourcePropertiesImpl) ServiceResourceProperties() BaseServiceResourcePropertiesImpl {
	return s.serviceResourceProperties
}

func (s RawServiceResourcePropertiesImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalServiceResourcePropertiesImplementation(input []byte) (ServiceResourceProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ServiceResourceProperties into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["serviceType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "DataTransfer") {
		var out DataTransferServiceResourceProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DataTransferServiceResourceProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GraphAPICompute") {
		var out GraphAPIComputeServiceResourceProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GraphAPIComputeServiceResourceProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MaterializedViewsBuilder") {
		var out MaterializedViewsBuilderServiceResourceProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MaterializedViewsBuilderServiceResourceProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SqlDedicatedGateway") {
		var out SqlDedicatedGatewayServiceResourceProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SqlDedicatedGatewayServiceResourceProperties: %+v", err)
		}
		return out, nil
	}

	var parent BaseServiceResourcePropertiesImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseServiceResourcePropertiesImpl: %+v", err)
	}

	return RawServiceResourcePropertiesImpl{
		serviceResourceProperties: parent,
		Type:                      value,
		Values:                    temp,
	}, nil

}
