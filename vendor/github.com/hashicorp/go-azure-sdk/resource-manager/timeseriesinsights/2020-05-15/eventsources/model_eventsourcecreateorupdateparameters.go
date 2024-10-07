package eventsources

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventSourceCreateOrUpdateParameters interface {
	EventSourceCreateOrUpdateParameters() BaseEventSourceCreateOrUpdateParametersImpl
}

var _ EventSourceCreateOrUpdateParameters = BaseEventSourceCreateOrUpdateParametersImpl{}

type BaseEventSourceCreateOrUpdateParametersImpl struct {
	Kind           EventSourceKind    `json:"kind"`
	LocalTimestamp *LocalTimestamp    `json:"localTimestamp,omitempty"`
	Location       string             `json:"location"`
	Tags           *map[string]string `json:"tags,omitempty"`
}

func (s BaseEventSourceCreateOrUpdateParametersImpl) EventSourceCreateOrUpdateParameters() BaseEventSourceCreateOrUpdateParametersImpl {
	return s
}

var _ EventSourceCreateOrUpdateParameters = RawEventSourceCreateOrUpdateParametersImpl{}

// RawEventSourceCreateOrUpdateParametersImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawEventSourceCreateOrUpdateParametersImpl struct {
	eventSourceCreateOrUpdateParameters BaseEventSourceCreateOrUpdateParametersImpl
	Type                                string
	Values                              map[string]interface{}
}

func (s RawEventSourceCreateOrUpdateParametersImpl) EventSourceCreateOrUpdateParameters() BaseEventSourceCreateOrUpdateParametersImpl {
	return s.eventSourceCreateOrUpdateParameters
}

func UnmarshalEventSourceCreateOrUpdateParametersImplementation(input []byte) (EventSourceCreateOrUpdateParameters, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling EventSourceCreateOrUpdateParameters into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["kind"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Microsoft.EventHub") {
		var out EventHubEventSourceCreateOrUpdateParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EventHubEventSourceCreateOrUpdateParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.IoTHub") {
		var out IoTHubEventSourceCreateOrUpdateParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into IoTHubEventSourceCreateOrUpdateParameters: %+v", err)
		}
		return out, nil
	}

	var parent BaseEventSourceCreateOrUpdateParametersImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseEventSourceCreateOrUpdateParametersImpl: %+v", err)
	}

	return RawEventSourceCreateOrUpdateParametersImpl{
		eventSourceCreateOrUpdateParameters: parent,
		Type:                                value,
		Values:                              temp,
	}, nil

}
