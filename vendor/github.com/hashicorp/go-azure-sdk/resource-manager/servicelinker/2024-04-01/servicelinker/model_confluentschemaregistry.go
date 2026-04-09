package servicelinker

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TargetServiceBase = ConfluentSchemaRegistry{}

type ConfluentSchemaRegistry struct {
	Endpoint *string `json:"endpoint,omitempty"`

	// Fields inherited from TargetServiceBase

	Type TargetServiceType `json:"type"`
}

func (s ConfluentSchemaRegistry) TargetServiceBase() BaseTargetServiceBaseImpl {
	return BaseTargetServiceBaseImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = ConfluentSchemaRegistry{}

func (s ConfluentSchemaRegistry) MarshalJSON() ([]byte, error) {
	type wrapper ConfluentSchemaRegistry
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ConfluentSchemaRegistry: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ConfluentSchemaRegistry: %+v", err)
	}

	decoded["type"] = "ConfluentSchemaRegistry"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ConfluentSchemaRegistry: %+v", err)
	}

	return encoded, nil
}
