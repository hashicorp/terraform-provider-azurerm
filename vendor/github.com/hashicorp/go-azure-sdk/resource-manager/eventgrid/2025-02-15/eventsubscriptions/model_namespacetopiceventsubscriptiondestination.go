package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EventSubscriptionDestination = NamespaceTopicEventSubscriptionDestination{}

type NamespaceTopicEventSubscriptionDestination struct {
	Properties *NamespaceTopicEventSubscriptionDestinationProperties `json:"properties,omitempty"`

	// Fields inherited from EventSubscriptionDestination

	EndpointType EndpointType `json:"endpointType"`
}

func (s NamespaceTopicEventSubscriptionDestination) EventSubscriptionDestination() BaseEventSubscriptionDestinationImpl {
	return BaseEventSubscriptionDestinationImpl{
		EndpointType: s.EndpointType,
	}
}

var _ json.Marshaler = NamespaceTopicEventSubscriptionDestination{}

func (s NamespaceTopicEventSubscriptionDestination) MarshalJSON() ([]byte, error) {
	type wrapper NamespaceTopicEventSubscriptionDestination
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling NamespaceTopicEventSubscriptionDestination: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling NamespaceTopicEventSubscriptionDestination: %+v", err)
	}

	decoded["endpointType"] = "NamespaceTopic"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling NamespaceTopicEventSubscriptionDestination: %+v", err)
	}

	return encoded, nil
}
