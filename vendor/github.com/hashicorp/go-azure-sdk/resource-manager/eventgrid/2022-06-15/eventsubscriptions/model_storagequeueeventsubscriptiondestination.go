package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EventSubscriptionDestination = StorageQueueEventSubscriptionDestination{}

type StorageQueueEventSubscriptionDestination struct {
	Properties *StorageQueueEventSubscriptionDestinationProperties `json:"properties,omitempty"`

	// Fields inherited from EventSubscriptionDestination

	EndpointType EndpointType `json:"endpointType"`
}

func (s StorageQueueEventSubscriptionDestination) EventSubscriptionDestination() BaseEventSubscriptionDestinationImpl {
	return BaseEventSubscriptionDestinationImpl{
		EndpointType: s.EndpointType,
	}
}

var _ json.Marshaler = StorageQueueEventSubscriptionDestination{}

func (s StorageQueueEventSubscriptionDestination) MarshalJSON() ([]byte, error) {
	type wrapper StorageQueueEventSubscriptionDestination
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling StorageQueueEventSubscriptionDestination: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling StorageQueueEventSubscriptionDestination: %+v", err)
	}

	decoded["endpointType"] = "StorageQueue"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling StorageQueueEventSubscriptionDestination: %+v", err)
	}

	return encoded, nil
}
