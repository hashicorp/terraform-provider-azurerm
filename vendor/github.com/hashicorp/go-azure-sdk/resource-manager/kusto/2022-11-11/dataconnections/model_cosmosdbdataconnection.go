package dataconnections

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DataConnection = CosmosDbDataConnection{}

type CosmosDbDataConnection struct {
	Properties *CosmosDbDataConnectionProperties `json:"properties,omitempty"`

	// Fields inherited from DataConnection
	Id       *string `json:"id,omitempty"`
	Location *string `json:"location,omitempty"`
	Name     *string `json:"name,omitempty"`
	Type     *string `json:"type,omitempty"`
}

var _ json.Marshaler = CosmosDbDataConnection{}

func (s CosmosDbDataConnection) MarshalJSON() ([]byte, error) {
	type wrapper CosmosDbDataConnection
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CosmosDbDataConnection: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CosmosDbDataConnection: %+v", err)
	}
	decoded["kind"] = "CosmosDb"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CosmosDbDataConnection: %+v", err)
	}

	return encoded, nil
}
