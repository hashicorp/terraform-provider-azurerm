package dataflows

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataFlow interface {
	DataFlow() BaseDataFlowImpl
}

var _ DataFlow = BaseDataFlowImpl{}

type BaseDataFlowImpl struct {
	Annotations *[]interface{}  `json:"annotations,omitempty"`
	Description *string         `json:"description,omitempty"`
	Folder      *DataFlowFolder `json:"folder,omitempty"`
	Type        string          `json:"type"`
}

func (s BaseDataFlowImpl) DataFlow() BaseDataFlowImpl {
	return s
}

var _ DataFlow = RawDataFlowImpl{}

// RawDataFlowImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawDataFlowImpl struct {
	dataFlow BaseDataFlowImpl
	Type     string
	Values   map[string]interface{}
}

func (s RawDataFlowImpl) DataFlow() BaseDataFlowImpl {
	return s.dataFlow
}

func UnmarshalDataFlowImplementation(input []byte) (DataFlow, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DataFlow into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Flowlet") {
		var out Flowlet
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into Flowlet: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MappingDataFlow") {
		var out MappingDataFlow
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MappingDataFlow: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "WranglingDataFlow") {
		var out WranglingDataFlow
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into WranglingDataFlow: %+v", err)
		}
		return out, nil
	}

	var parent BaseDataFlowImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseDataFlowImpl: %+v", err)
	}

	return RawDataFlowImpl{
		dataFlow: parent,
		Type:     value,
		Values:   temp,
	}, nil

}
