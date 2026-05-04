package dataconnections

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataConnection interface {
	DataConnection() BaseDataConnectionImpl
}

var _ DataConnection = BaseDataConnectionImpl{}

type BaseDataConnectionImpl struct {
	Id       *string            `json:"id,omitempty"`
	Kind     DataConnectionKind `json:"kind"`
	Location *string            `json:"location,omitempty"`
	Name     *string            `json:"name,omitempty"`
	Type     *string            `json:"type,omitempty"`
}

func (s BaseDataConnectionImpl) DataConnection() BaseDataConnectionImpl {
	return s
}

var _ DataConnection = RawDataConnectionImpl{}

// RawDataConnectionImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawDataConnectionImpl struct {
	dataConnection BaseDataConnectionImpl
	Type           string
	Values         map[string]interface{}
}

func (s RawDataConnectionImpl) DataConnection() BaseDataConnectionImpl {
	return s.dataConnection
}

func UnmarshalDataConnectionImplementation(input []byte) (DataConnection, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DataConnection into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["kind"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "CosmosDb") {
		var out CosmosDbDataConnection
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CosmosDbDataConnection: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "EventGrid") {
		var out EventGridDataConnection
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EventGridDataConnection: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "EventHub") {
		var out EventHubDataConnection
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EventHubDataConnection: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "IotHub") {
		var out IotHubDataConnection
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into IotHubDataConnection: %+v", err)
		}
		return out, nil
	}

	var parent BaseDataConnectionImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseDataConnectionImpl: %+v", err)
	}

	return RawDataConnectionImpl{
		dataConnection: parent,
		Type:           value,
		Values:         temp,
	}, nil

}
