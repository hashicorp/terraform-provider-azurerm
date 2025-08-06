package inputs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ StreamInputDataSource = EventHubStreamInputDataSource{}

type EventHubStreamInputDataSource struct {
	Properties *EventHubStreamInputDataSourceProperties `json:"properties,omitempty"`

	// Fields inherited from StreamInputDataSource

	Type string `json:"type"`
}

func (s EventHubStreamInputDataSource) StreamInputDataSource() BaseStreamInputDataSourceImpl {
	return BaseStreamInputDataSourceImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = EventHubStreamInputDataSource{}

func (s EventHubStreamInputDataSource) MarshalJSON() ([]byte, error) {
	type wrapper EventHubStreamInputDataSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling EventHubStreamInputDataSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling EventHubStreamInputDataSource: %+v", err)
	}

	decoded["type"] = "Microsoft.ServiceBus/EventHub"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling EventHubStreamInputDataSource: %+v", err)
	}

	return encoded, nil
}
