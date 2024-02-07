package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EventSubscriptionDestination = ServiceBusTopicEventSubscriptionDestination{}

type ServiceBusTopicEventSubscriptionDestination struct {
	Properties *ServiceBusTopicEventSubscriptionDestinationProperties `json:"properties,omitempty"`

	// Fields inherited from EventSubscriptionDestination
}

var _ json.Marshaler = ServiceBusTopicEventSubscriptionDestination{}

func (s ServiceBusTopicEventSubscriptionDestination) MarshalJSON() ([]byte, error) {
	type wrapper ServiceBusTopicEventSubscriptionDestination
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ServiceBusTopicEventSubscriptionDestination: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ServiceBusTopicEventSubscriptionDestination: %+v", err)
	}
	decoded["endpointType"] = "ServiceBusTopic"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ServiceBusTopicEventSubscriptionDestination: %+v", err)
	}

	return encoded, nil
}
