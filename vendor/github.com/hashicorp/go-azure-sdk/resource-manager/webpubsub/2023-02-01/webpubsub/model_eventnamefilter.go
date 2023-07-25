package webpubsub

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EventListenerFilter = EventNameFilter{}

type EventNameFilter struct {
	SystemEvents     *[]string `json:"systemEvents,omitempty"`
	UserEventPattern *string   `json:"userEventPattern,omitempty"`

	// Fields inherited from EventListenerFilter
}

var _ json.Marshaler = EventNameFilter{}

func (s EventNameFilter) MarshalJSON() ([]byte, error) {
	type wrapper EventNameFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling EventNameFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling EventNameFilter: %+v", err)
	}
	decoded["type"] = "EventName"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling EventNameFilter: %+v", err)
	}

	return encoded, nil
}
