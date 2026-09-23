package eventsubscriptions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeadLetterDestination interface {
	DeadLetterDestination() BaseDeadLetterDestinationImpl
}

var _ DeadLetterDestination = BaseDeadLetterDestinationImpl{}

type BaseDeadLetterDestinationImpl struct {
	EndpointType DeadLetterEndPointType `json:"endpointType"`
}

func (s BaseDeadLetterDestinationImpl) DeadLetterDestination() BaseDeadLetterDestinationImpl {
	return s
}

var _ DeadLetterDestination = RawDeadLetterDestinationImpl{}

// RawDeadLetterDestinationImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawDeadLetterDestinationImpl struct {
	deadLetterDestination BaseDeadLetterDestinationImpl
	Type                  string
	Values                map[string]interface{}
}

func (s RawDeadLetterDestinationImpl) DeadLetterDestination() BaseDeadLetterDestinationImpl {
	return s.deadLetterDestination
}

func (s RawDeadLetterDestinationImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalDeadLetterDestinationImplementation(input []byte) (DeadLetterDestination, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DeadLetterDestination into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["endpointType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "StorageBlob") {
		var out StorageBlobDeadLetterDestination
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StorageBlobDeadLetterDestination: %+v", err)
		}
		return out, nil
	}

	var parent BaseDeadLetterDestinationImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseDeadLetterDestinationImpl: %+v", err)
	}

	return RawDeadLetterDestinationImpl{
		deadLetterDestination: parent,
		Type:                  value,
		Values:                temp,
	}, nil

}
