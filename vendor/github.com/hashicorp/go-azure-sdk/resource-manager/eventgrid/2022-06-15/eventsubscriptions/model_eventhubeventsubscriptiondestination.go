package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EventSubscriptionDestination = EventHubEventSubscriptionDestination{}

type EventHubEventSubscriptionDestination struct {
	Properties *EventHubEventSubscriptionDestinationProperties `json:"properties,omitempty"`

	// Fields inherited from EventSubscriptionDestination
}

var _ json.Marshaler = EventHubEventSubscriptionDestination{}

func (s EventHubEventSubscriptionDestination) MarshalJSON() ([]byte, error) {
	type wrapper EventHubEventSubscriptionDestination
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling EventHubEventSubscriptionDestination: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling EventHubEventSubscriptionDestination: %+v", err)
	}
	decoded["endpointType"] = "EventHub"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling EventHubEventSubscriptionDestination: %+v", err)
	}

	return encoded, nil
}
