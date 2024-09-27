package eventsources

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventSourceResource interface {
	EventSourceResource() BaseEventSourceResourceImpl
}

var _ EventSourceResource = BaseEventSourceResourceImpl{}

type BaseEventSourceResourceImpl struct {
	Id       *string            `json:"id,omitempty"`
	Kind     Kind               `json:"kind"`
	Location string             `json:"location"`
	Name     *string            `json:"name,omitempty"`
	Tags     *map[string]string `json:"tags,omitempty"`
	Type     *string            `json:"type,omitempty"`
}

func (s BaseEventSourceResourceImpl) EventSourceResource() BaseEventSourceResourceImpl {
	return s
}

var _ EventSourceResource = RawEventSourceResourceImpl{}

// RawEventSourceResourceImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawEventSourceResourceImpl struct {
	eventSourceResource BaseEventSourceResourceImpl
	Type                string
	Values              map[string]interface{}
}

func (s RawEventSourceResourceImpl) EventSourceResource() BaseEventSourceResourceImpl {
	return s.eventSourceResource
}

func UnmarshalEventSourceResourceImplementation(input []byte) (EventSourceResource, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling EventSourceResource into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["kind"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Microsoft.EventHub") {
		var out EventHubEventSourceResource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EventHubEventSourceResource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Microsoft.IoTHub") {
		var out IoTHubEventSourceResource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into IoTHubEventSourceResource: %+v", err)
		}
		return out, nil
	}

	var parent BaseEventSourceResourceImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseEventSourceResourceImpl: %+v", err)
	}

	return RawEventSourceResourceImpl{
		eventSourceResource: parent,
		Type:                value,
		Values:              temp,
	}, nil

}
