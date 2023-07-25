package eventsources

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EventSourceCreateOrUpdateParameters = EventHubEventSourceCreateOrUpdateParameters{}

type EventHubEventSourceCreateOrUpdateParameters struct {
	Properties EventHubEventSourceCreationProperties `json:"properties"`

	// Fields inherited from EventSourceCreateOrUpdateParameters
	LocalTimestamp *LocalTimestamp    `json:"localTimestamp,omitempty"`
	Location       string             `json:"location"`
	Tags           *map[string]string `json:"tags,omitempty"`
}

var _ json.Marshaler = EventHubEventSourceCreateOrUpdateParameters{}

func (s EventHubEventSourceCreateOrUpdateParameters) MarshalJSON() ([]byte, error) {
	type wrapper EventHubEventSourceCreateOrUpdateParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling EventHubEventSourceCreateOrUpdateParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling EventHubEventSourceCreateOrUpdateParameters: %+v", err)
	}
	decoded["kind"] = "Microsoft.EventHub"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling EventHubEventSourceCreateOrUpdateParameters: %+v", err)
	}

	return encoded, nil
}
