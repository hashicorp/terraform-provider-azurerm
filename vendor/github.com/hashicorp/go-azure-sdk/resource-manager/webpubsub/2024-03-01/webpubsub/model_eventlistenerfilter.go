package webpubsub

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventListenerFilter interface {
	EventListenerFilter() BaseEventListenerFilterImpl
}

var _ EventListenerFilter = BaseEventListenerFilterImpl{}

type BaseEventListenerFilterImpl struct {
	Type EventListenerFilterDiscriminator `json:"type"`
}

func (s BaseEventListenerFilterImpl) EventListenerFilter() BaseEventListenerFilterImpl {
	return s
}

var _ EventListenerFilter = RawEventListenerFilterImpl{}

// RawEventListenerFilterImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawEventListenerFilterImpl struct {
	eventListenerFilter BaseEventListenerFilterImpl
	Type                string
	Values              map[string]interface{}
}

func (s RawEventListenerFilterImpl) EventListenerFilter() BaseEventListenerFilterImpl {
	return s.eventListenerFilter
}

func (s RawEventListenerFilterImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalEventListenerFilterImplementation(input []byte) (EventListenerFilter, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling EventListenerFilter into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "EventName") {
		var out EventNameFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EventNameFilter: %+v", err)
		}
		return out, nil
	}

	var parent BaseEventListenerFilterImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseEventListenerFilterImpl: %+v", err)
	}

	return RawEventListenerFilterImpl{
		eventListenerFilter: parent,
		Type:                value,
		Values:              temp,
	}, nil

}
