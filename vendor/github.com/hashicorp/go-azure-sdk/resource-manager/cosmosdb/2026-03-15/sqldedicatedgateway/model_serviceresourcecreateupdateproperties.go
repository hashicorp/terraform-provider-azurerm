package sqldedicatedgateway

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceResourceCreateUpdateProperties interface {
	ServiceResourceCreateUpdateProperties() BaseServiceResourceCreateUpdatePropertiesImpl
}

var _ ServiceResourceCreateUpdateProperties = BaseServiceResourceCreateUpdatePropertiesImpl{}

type BaseServiceResourceCreateUpdatePropertiesImpl struct {
	InstanceCount *int64       `json:"instanceCount,omitempty"`
	InstanceSize  *ServiceSize `json:"instanceSize,omitempty"`
	ServiceType   ServiceType  `json:"serviceType"`
}

func (s BaseServiceResourceCreateUpdatePropertiesImpl) ServiceResourceCreateUpdateProperties() BaseServiceResourceCreateUpdatePropertiesImpl {
	return s
}

var _ ServiceResourceCreateUpdateProperties = RawServiceResourceCreateUpdatePropertiesImpl{}

// RawServiceResourceCreateUpdatePropertiesImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawServiceResourceCreateUpdatePropertiesImpl struct {
	serviceResourceCreateUpdateProperties BaseServiceResourceCreateUpdatePropertiesImpl
	Type                                  string
	Values                                map[string]interface{}
}

func (s RawServiceResourceCreateUpdatePropertiesImpl) ServiceResourceCreateUpdateProperties() BaseServiceResourceCreateUpdatePropertiesImpl {
	return s.serviceResourceCreateUpdateProperties
}

func (s RawServiceResourceCreateUpdatePropertiesImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalServiceResourceCreateUpdatePropertiesImplementation(input []byte) (ServiceResourceCreateUpdateProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ServiceResourceCreateUpdateProperties into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["serviceType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "DataTransfer") {
		var out DataTransferServiceResourceCreateUpdateProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DataTransferServiceResourceCreateUpdateProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GraphAPICompute") {
		var out GraphAPIComputeServiceResourceCreateUpdateProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GraphAPIComputeServiceResourceCreateUpdateProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MaterializedViewsBuilder") {
		var out MaterializedViewsBuilderServiceResourceCreateUpdateProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MaterializedViewsBuilderServiceResourceCreateUpdateProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SqlDedicatedGateway") {
		var out SqlDedicatedGatewayServiceResourceCreateUpdateProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SqlDedicatedGatewayServiceResourceCreateUpdateProperties: %+v", err)
		}
		return out, nil
	}

	var parent BaseServiceResourceCreateUpdatePropertiesImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseServiceResourceCreateUpdatePropertiesImpl: %+v", err)
	}

	return RawServiceResourceCreateUpdatePropertiesImpl{
		serviceResourceCreateUpdateProperties: parent,
		Type:                                  value,
		Values:                                temp,
	}, nil

}
