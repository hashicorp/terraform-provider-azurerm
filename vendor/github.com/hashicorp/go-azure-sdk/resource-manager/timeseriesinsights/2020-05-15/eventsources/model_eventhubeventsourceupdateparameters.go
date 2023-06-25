package eventsources

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EventSourceUpdateParameters = EventHubEventSourceUpdateParameters{}

type EventHubEventSourceUpdateParameters struct {
	Properties *EventHubEventSourceMutableProperties `json:"properties,omitempty"`

	// Fields inherited from EventSourceUpdateParameters
	Tags *map[string]string `json:"tags,omitempty"`
}

var _ json.Marshaler = EventHubEventSourceUpdateParameters{}

func (s EventHubEventSourceUpdateParameters) MarshalJSON() ([]byte, error) {
	type wrapper EventHubEventSourceUpdateParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling EventHubEventSourceUpdateParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling EventHubEventSourceUpdateParameters: %+v", err)
	}
	decoded["kind"] = "Microsoft.EventHub"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling EventHubEventSourceUpdateParameters: %+v", err)
	}

	return encoded, nil
}
