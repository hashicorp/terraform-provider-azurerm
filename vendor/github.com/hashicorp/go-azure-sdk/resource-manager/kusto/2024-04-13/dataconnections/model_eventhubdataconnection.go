package dataconnections

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DataConnection = EventHubDataConnection{}

type EventHubDataConnection struct {
	Properties *EventHubConnectionProperties `json:"properties,omitempty"`

	// Fields inherited from DataConnection

	Id       *string            `json:"id,omitempty"`
	Kind     DataConnectionKind `json:"kind"`
	Location *string            `json:"location,omitempty"`
	Name     *string            `json:"name,omitempty"`
	Type     *string            `json:"type,omitempty"`
}

func (s EventHubDataConnection) DataConnection() BaseDataConnectionImpl {
	return BaseDataConnectionImpl{
		Id:       s.Id,
		Kind:     s.Kind,
		Location: s.Location,
		Name:     s.Name,
		Type:     s.Type,
	}
}

var _ json.Marshaler = EventHubDataConnection{}

func (s EventHubDataConnection) MarshalJSON() ([]byte, error) {
	type wrapper EventHubDataConnection
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling EventHubDataConnection: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling EventHubDataConnection: %+v", err)
	}

	decoded["kind"] = "EventHub"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling EventHubDataConnection: %+v", err)
	}

	return encoded, nil
}
