package eventsources

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EventSourceUpdateParameters = IoTHubEventSourceUpdateParameters{}

type IoTHubEventSourceUpdateParameters struct {
	Properties *IoTHubEventSourceMutableProperties `json:"properties,omitempty"`

	// Fields inherited from EventSourceUpdateParameters
	Tags *map[string]string `json:"tags,omitempty"`
}

var _ json.Marshaler = IoTHubEventSourceUpdateParameters{}

func (s IoTHubEventSourceUpdateParameters) MarshalJSON() ([]byte, error) {
	type wrapper IoTHubEventSourceUpdateParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling IoTHubEventSourceUpdateParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling IoTHubEventSourceUpdateParameters: %+v", err)
	}
	decoded["kind"] = "Microsoft.IoTHub"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling IoTHubEventSourceUpdateParameters: %+v", err)
	}

	return encoded, nil
}
