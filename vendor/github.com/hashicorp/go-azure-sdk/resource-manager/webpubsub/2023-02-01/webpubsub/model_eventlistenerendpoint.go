package webpubsub

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventListenerEndpoint interface {
	EventListenerEndpoint() BaseEventListenerEndpointImpl
}

var _ EventListenerEndpoint = BaseEventListenerEndpointImpl{}

type BaseEventListenerEndpointImpl struct {
	Type EventListenerEndpointDiscriminator `json:"type"`
}

func (s BaseEventListenerEndpointImpl) EventListenerEndpoint() BaseEventListenerEndpointImpl {
	return s
}

var _ EventListenerEndpoint = RawEventListenerEndpointImpl{}

// RawEventListenerEndpointImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawEventListenerEndpointImpl struct {
	eventListenerEndpoint BaseEventListenerEndpointImpl
	Type                  string
	Values                map[string]interface{}
}

func (s RawEventListenerEndpointImpl) EventListenerEndpoint() BaseEventListenerEndpointImpl {
	return s.eventListenerEndpoint
}

func UnmarshalEventListenerEndpointImplementation(input []byte) (EventListenerEndpoint, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling EventListenerEndpoint into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "EventHub") {
		var out EventHubEndpoint
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EventHubEndpoint: %+v", err)
		}
		return out, nil
	}

	var parent BaseEventListenerEndpointImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseEventListenerEndpointImpl: %+v", err)
	}

	return RawEventListenerEndpointImpl{
		eventListenerEndpoint: parent,
		Type:                  value,
		Values:                temp,
	}, nil

}
