package eventsources

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventSourceUpdateParameters interface {
	EventSourceUpdateParameters() BaseEventSourceUpdateParametersImpl
}

var _ EventSourceUpdateParameters = BaseEventSourceUpdateParametersImpl{}

type BaseEventSourceUpdateParametersImpl struct {
	Kind EventSourceKind    `json:"kind"`
	Tags *map[string]string `json:"tags,omitempty"`
}

func (s BaseEventSourceUpdateParametersImpl) EventSourceUpdateParameters() BaseEventSourceUpdateParametersImpl {
	return s
}

var _ EventSourceUpdateParameters = RawEventSourceUpdateParametersImpl{}

// RawEventSourceUpdateParametersImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawEventSourceUpdateParametersImpl struct {
	eventSourceUpdateParameters BaseEventSourceUpdateParametersImpl
	Type                        string
	Values                      map[string]interface{}
}

func (s RawEventSourceUpdateParametersImpl) EventSourceUpdateParameters() BaseEventSourceUpdateParametersImpl {
	return s.eventSourceUpdateParameters
}

func UnmarshalEventSourceUpdateParametersImplementation(input []byte) (EventSourceUpdateParameters, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling EventSourceUpdateParameters into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["kind"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Microsoft.EventHub") {
		var out EventHubEventSourceUpdateParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EventHubEventSourceUpdateParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.IoTHub") {
		var out IoTHubEventSourceUpdateParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into IoTHubEventSourceUpdateParameters: %+v", err)
		}
		return out, nil
	}

	var parent BaseEventSourceUpdateParametersImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseEventSourceUpdateParametersImpl: %+v", err)
	}

	return RawEventSourceUpdateParametersImpl{
		eventSourceUpdateParameters: parent,
		Type:                        value,
		Values:                      temp,
	}, nil

}
