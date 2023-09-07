package dataconnections

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataConnection interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawDataConnectionImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalDataConnectionImplementation(input []byte) (DataConnection, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DataConnection into map[string]interface: %+v", err)
	}

	value, ok := temp["kind"].(string)
	if !ok {
		return nil, nil
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

	out := RawDataConnectionImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
