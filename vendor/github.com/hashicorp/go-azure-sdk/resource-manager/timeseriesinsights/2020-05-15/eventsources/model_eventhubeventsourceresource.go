package eventsources

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EventSourceResource = EventHubEventSourceResource{}

type EventHubEventSourceResource struct {
	Properties EventHubEventSourceCommonProperties `json:"properties"`

	// Fields inherited from EventSourceResource
	Id       *string            `json:"id,omitempty"`
	Location string             `json:"location"`
	Name     *string            `json:"name,omitempty"`
	Tags     *map[string]string `json:"tags,omitempty"`
	Type     *string            `json:"type,omitempty"`
}

var _ json.Marshaler = EventHubEventSourceResource{}

func (s EventHubEventSourceResource) MarshalJSON() ([]byte, error) {
	type wrapper EventHubEventSourceResource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling EventHubEventSourceResource: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling EventHubEventSourceResource: %+v", err)
	}
	decoded["kind"] = "Microsoft.EventHub"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling EventHubEventSourceResource: %+v", err)
	}

	return encoded, nil
}
