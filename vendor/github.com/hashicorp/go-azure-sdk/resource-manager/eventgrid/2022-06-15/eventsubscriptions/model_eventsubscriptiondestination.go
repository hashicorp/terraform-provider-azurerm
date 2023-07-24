package eventsubscriptions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventSubscriptionDestination interface {
}

func unmarshalEventSubscriptionDestinationImplementation(input []byte) (EventSubscriptionDestination, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling EventSubscriptionDestination into map[string]interface: %+v", err)
	}

	value, ok := temp["endpointType"].(string)
	if !ok {
		return nil, nil
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

	type RawEventSubscriptionDestinationImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawEventSubscriptionDestinationImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
