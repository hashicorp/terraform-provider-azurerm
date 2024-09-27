package endpoints

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DigitalTwinsEndpointResourceProperties interface {
	DigitalTwinsEndpointResourceProperties() BaseDigitalTwinsEndpointResourcePropertiesImpl
}

var _ DigitalTwinsEndpointResourceProperties = BaseDigitalTwinsEndpointResourcePropertiesImpl{}

type BaseDigitalTwinsEndpointResourcePropertiesImpl struct {
	AuthenticationType *AuthenticationType        `json:"authenticationType,omitempty"`
	CreatedTime        *string                    `json:"createdTime,omitempty"`
	DeadLetterSecret   *string                    `json:"deadLetterSecret,omitempty"`
	DeadLetterUri      *string                    `json:"deadLetterUri,omitempty"`
	EndpointType       EndpointType               `json:"endpointType"`
	Identity           *ManagedIdentityReference  `json:"identity,omitempty"`
	ProvisioningState  *EndpointProvisioningState `json:"provisioningState,omitempty"`
}

func (s BaseDigitalTwinsEndpointResourcePropertiesImpl) DigitalTwinsEndpointResourceProperties() BaseDigitalTwinsEndpointResourcePropertiesImpl {
	return s
}

var _ DigitalTwinsEndpointResourceProperties = RawDigitalTwinsEndpointResourcePropertiesImpl{}

// RawDigitalTwinsEndpointResourcePropertiesImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawDigitalTwinsEndpointResourcePropertiesImpl struct {
	digitalTwinsEndpointResourceProperties BaseDigitalTwinsEndpointResourcePropertiesImpl
	Type                                   string
	Values                                 map[string]interface{}
}

func (s RawDigitalTwinsEndpointResourcePropertiesImpl) DigitalTwinsEndpointResourceProperties() BaseDigitalTwinsEndpointResourcePropertiesImpl {
	return s.digitalTwinsEndpointResourceProperties
}

func UnmarshalDigitalTwinsEndpointResourcePropertiesImplementation(input []byte) (DigitalTwinsEndpointResourceProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DigitalTwinsEndpointResourceProperties into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["endpointType"]; ok {
		value = fmt.Sprintf("%v", v)
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

	var parent BaseDigitalTwinsEndpointResourcePropertiesImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseDigitalTwinsEndpointResourcePropertiesImpl: %+v", err)
	}

	return RawDigitalTwinsEndpointResourcePropertiesImpl{
		digitalTwinsEndpointResourceProperties: parent,
		Type:                                   value,
		Values:                                 temp,
	}, nil

}
