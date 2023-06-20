package inputs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ StreamInputDataSource = EventHubV2StreamInputDataSource{}

type EventHubV2StreamInputDataSource struct {
	Properties *EventHubStreamInputDataSourceProperties `json:"properties,omitempty"`

	// Fields inherited from StreamInputDataSource
}

var _ json.Marshaler = EventHubV2StreamInputDataSource{}

func (s EventHubV2StreamInputDataSource) MarshalJSON() ([]byte, error) {
	type wrapper EventHubV2StreamInputDataSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling EventHubV2StreamInputDataSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling EventHubV2StreamInputDataSource: %+v", err)
	}
	decoded["type"] = "Microsoft.EventHub/EventHub"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling EventHubV2StreamInputDataSource: %+v", err)
	}

	return encoded, nil
}
