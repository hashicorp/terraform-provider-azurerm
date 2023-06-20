package eventsources

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EventSourceCreateOrUpdateParameters = IoTHubEventSourceCreateOrUpdateParameters{}

type IoTHubEventSourceCreateOrUpdateParameters struct {
	Properties IoTHubEventSourceCreationProperties `json:"properties"`

	// Fields inherited from EventSourceCreateOrUpdateParameters
	LocalTimestamp *LocalTimestamp    `json:"localTimestamp,omitempty"`
	Location       string             `json:"location"`
	Tags           *map[string]string `json:"tags,omitempty"`
}

var _ json.Marshaler = IoTHubEventSourceCreateOrUpdateParameters{}

func (s IoTHubEventSourceCreateOrUpdateParameters) MarshalJSON() ([]byte, error) {
	type wrapper IoTHubEventSourceCreateOrUpdateParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling IoTHubEventSourceCreateOrUpdateParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling IoTHubEventSourceCreateOrUpdateParameters: %+v", err)
	}
	decoded["kind"] = "Microsoft.IoTHub"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling IoTHubEventSourceCreateOrUpdateParameters: %+v", err)
	}

	return encoded, nil
}
