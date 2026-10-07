package eventsubscriptions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeliveryAttributeMapping interface {
	DeliveryAttributeMapping() BaseDeliveryAttributeMappingImpl
}

var _ DeliveryAttributeMapping = BaseDeliveryAttributeMappingImpl{}

type BaseDeliveryAttributeMappingImpl struct {
	Name *string                      `json:"name,omitempty"`
	Type DeliveryAttributeMappingType `json:"type"`
}

func (s BaseDeliveryAttributeMappingImpl) DeliveryAttributeMapping() BaseDeliveryAttributeMappingImpl {
	return s
}

var _ DeliveryAttributeMapping = RawDeliveryAttributeMappingImpl{}

// RawDeliveryAttributeMappingImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawDeliveryAttributeMappingImpl struct {
	deliveryAttributeMapping BaseDeliveryAttributeMappingImpl
	Type                     string
	Values                   map[string]interface{}
}

func (s RawDeliveryAttributeMappingImpl) DeliveryAttributeMapping() BaseDeliveryAttributeMappingImpl {
	return s.deliveryAttributeMapping
}

func (s RawDeliveryAttributeMappingImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalDeliveryAttributeMappingImplementation(input []byte) (DeliveryAttributeMapping, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryAttributeMapping into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Dynamic") {
		var out DynamicDeliveryAttributeMapping
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DynamicDeliveryAttributeMapping: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Static") {
		var out StaticDeliveryAttributeMapping
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StaticDeliveryAttributeMapping: %+v", err)
		}
		return out, nil
	}

	var parent BaseDeliveryAttributeMappingImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseDeliveryAttributeMappingImpl: %+v", err)
	}

	return RawDeliveryAttributeMappingImpl{
		deliveryAttributeMapping: parent,
		Type:                     value,
		Values:                   temp,
	}, nil

}
