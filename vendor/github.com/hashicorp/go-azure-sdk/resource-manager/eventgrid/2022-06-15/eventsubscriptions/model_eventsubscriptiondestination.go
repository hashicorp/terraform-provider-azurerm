package eventsubscriptions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventSubscriptionDestination interface {
	EventSubscriptionDestination() BaseEventSubscriptionDestinationImpl
}

var _ EventSubscriptionDestination = BaseEventSubscriptionDestinationImpl{}

type BaseEventSubscriptionDestinationImpl struct {
	EndpointType EndpointType `json:"endpointType"`
}

func (s BaseEventSubscriptionDestinationImpl) EventSubscriptionDestination() BaseEventSubscriptionDestinationImpl {
	return s
}

var _ EventSubscriptionDestination = RawEventSubscriptionDestinationImpl{}

// RawEventSubscriptionDestinationImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawEventSubscriptionDestinationImpl struct {
	eventSubscriptionDestination BaseEventSubscriptionDestinationImpl
	Type                         string
	Values                       map[string]interface{}
}

func (s RawEventSubscriptionDestinationImpl) EventSubscriptionDestination() BaseEventSubscriptionDestinationImpl {
	return s.eventSubscriptionDestination
}

func UnmarshalEventSubscriptionDestinationImplementation(input []byte) (EventSubscriptionDestination, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling EventSubscriptionDestination into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["endpointType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AzureFunction") {
		var out AzureFunctionEventSubscriptionDestination
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureFunctionEventSubscriptionDestination: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "EventHub") {
		var out EventHubEventSubscriptionDestination
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EventHubEventSubscriptionDestination: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HybridConnection") {
		var out HybridConnectionEventSubscriptionDestination
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HybridConnectionEventSubscriptionDestination: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ServiceBusQueue") {
		var out ServiceBusQueueEventSubscriptionDestination
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServiceBusQueueEventSubscriptionDestination: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ServiceBusTopic") {
		var out ServiceBusTopicEventSubscriptionDestination
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServiceBusTopicEventSubscriptionDestination: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "StorageQueue") {
		var out StorageQueueEventSubscriptionDestination
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StorageQueueEventSubscriptionDestination: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "WebHook") {
		var out WebHookEventSubscriptionDestination
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into WebHookEventSubscriptionDestination: %+v", err)
		}
		return out, nil
	}

	var parent BaseEventSubscriptionDestinationImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseEventSubscriptionDestinationImpl: %+v", err)
	}

	return RawEventSubscriptionDestinationImpl{
		eventSubscriptionDestination: parent,
		Type:                         value,
		Values:                       temp,
	}, nil

}
