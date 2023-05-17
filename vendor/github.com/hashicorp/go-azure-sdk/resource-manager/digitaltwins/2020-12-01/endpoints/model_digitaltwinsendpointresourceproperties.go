package endpoints

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DigitalTwinsEndpointResourceProperties interface {
}

func unmarshalDigitalTwinsEndpointResourcePropertiesImplementation(input []byte) (DigitalTwinsEndpointResourceProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DigitalTwinsEndpointResourceProperties into map[string]interface: %+v", err)
	}

	value, ok := temp["endpointType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "EventGrid") {
		var out EventGrid
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EventGrid: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "EventHub") {
		var out EventHub
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EventHub: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ServiceBus") {
		var out ServiceBus
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServiceBus: %+v", err)
		}
		return out, nil
	}

	type RawDigitalTwinsEndpointResourcePropertiesImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawDigitalTwinsEndpointResourcePropertiesImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
