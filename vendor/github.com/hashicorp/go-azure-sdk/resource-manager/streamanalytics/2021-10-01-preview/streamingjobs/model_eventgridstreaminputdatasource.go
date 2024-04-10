package streamingjobs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ StreamInputDataSource = EventGridStreamInputDataSource{}

type EventGridStreamInputDataSource struct {
	Properties *EventGridStreamInputDataSourceProperties `json:"properties,omitempty"`

	// Fields inherited from StreamInputDataSource
}

var _ json.Marshaler = EventGridStreamInputDataSource{}

func (s EventGridStreamInputDataSource) MarshalJSON() ([]byte, error) {
	type wrapper EventGridStreamInputDataSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling EventGridStreamInputDataSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling EventGridStreamInputDataSource: %+v", err)
	}
	decoded["type"] = "Microsoft.EventGrid/EventSubscriptions"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling EventGridStreamInputDataSource: %+v", err)
	}

	return encoded, nil
}
