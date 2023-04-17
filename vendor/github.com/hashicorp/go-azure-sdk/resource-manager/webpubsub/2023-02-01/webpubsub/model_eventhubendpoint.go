package webpubsub

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EventListenerEndpoint = EventHubEndpoint{}

type EventHubEndpoint struct {
	EventHubName            string `json:"eventHubName"`
	FullyQualifiedNamespace string `json:"fullyQualifiedNamespace"`

	// Fields inherited from EventListenerEndpoint
}

var _ json.Marshaler = EventHubEndpoint{}

func (s EventHubEndpoint) MarshalJSON() ([]byte, error) {
	type wrapper EventHubEndpoint
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling EventHubEndpoint: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling EventHubEndpoint: %+v", err)
	}
	decoded["type"] = "EventHub"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling EventHubEndpoint: %+v", err)
	}

	return encoded, nil
}
