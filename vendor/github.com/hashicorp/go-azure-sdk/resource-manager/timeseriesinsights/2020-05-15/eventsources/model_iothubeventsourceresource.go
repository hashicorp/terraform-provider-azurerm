package eventsources

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EventSourceResource = IoTHubEventSourceResource{}

type IoTHubEventSourceResource struct {
	Properties IoTHubEventSourceCommonProperties `json:"properties"`

	// Fields inherited from EventSourceResource
	Id       *string            `json:"id,omitempty"`
	Location string             `json:"location"`
	Name     *string            `json:"name,omitempty"`
	Tags     *map[string]string `json:"tags,omitempty"`
	Type     *string            `json:"type,omitempty"`
}

var _ json.Marshaler = IoTHubEventSourceResource{}

func (s IoTHubEventSourceResource) MarshalJSON() ([]byte, error) {
	type wrapper IoTHubEventSourceResource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling IoTHubEventSourceResource: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling IoTHubEventSourceResource: %+v", err)
	}
	decoded["kind"] = "Microsoft.IoTHub"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling IoTHubEventSourceResource: %+v", err)
	}

	return encoded, nil
}
