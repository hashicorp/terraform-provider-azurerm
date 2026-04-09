package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EventSubscriptionDestination = ServiceBusQueueEventSubscriptionDestination{}

type ServiceBusQueueEventSubscriptionDestination struct {
	Properties *ServiceBusQueueEventSubscriptionDestinationProperties `json:"properties,omitempty"`

	// Fields inherited from EventSubscriptionDestination

	EndpointType EndpointType `json:"endpointType"`
}

func (s ServiceBusQueueEventSubscriptionDestination) EventSubscriptionDestination() BaseEventSubscriptionDestinationImpl {
	return BaseEventSubscriptionDestinationImpl{
		EndpointType: s.EndpointType,
	}
}

var _ json.Marshaler = ServiceBusQueueEventSubscriptionDestination{}

func (s ServiceBusQueueEventSubscriptionDestination) MarshalJSON() ([]byte, error) {
	type wrapper ServiceBusQueueEventSubscriptionDestination
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ServiceBusQueueEventSubscriptionDestination: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ServiceBusQueueEventSubscriptionDestination: %+v", err)
	}

	decoded["endpointType"] = "ServiceBusQueue"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ServiceBusQueueEventSubscriptionDestination: %+v", err)
	}

	return encoded, nil
}
