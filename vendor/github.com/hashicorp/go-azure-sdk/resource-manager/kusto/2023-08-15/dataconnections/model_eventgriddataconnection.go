package dataconnections

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DataConnection = EventGridDataConnection{}

type EventGridDataConnection struct {
	Properties *EventGridConnectionProperties `json:"properties,omitempty"`

	// Fields inherited from DataConnection
	Id       *string `json:"id,omitempty"`
	Location *string `json:"location,omitempty"`
	Name     *string `json:"name,omitempty"`
	Type     *string `json:"type,omitempty"`
}

var _ json.Marshaler = EventGridDataConnection{}

func (s EventGridDataConnection) MarshalJSON() ([]byte, error) {
	type wrapper EventGridDataConnection
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling EventGridDataConnection: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling EventGridDataConnection: %+v", err)
	}
	decoded["kind"] = "EventGrid"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling EventGridDataConnection: %+v", err)
	}

	return encoded, nil
}
