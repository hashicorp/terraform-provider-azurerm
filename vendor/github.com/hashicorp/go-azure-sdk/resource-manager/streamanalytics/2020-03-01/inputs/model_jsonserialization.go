package inputs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Serialization = JsonSerialization{}

type JsonSerialization struct {
	Properties *JsonSerializationProperties `json:"properties,omitempty"`

	// Fields inherited from Serialization

	Type EventSerializationType `json:"type"`
}

func (s JsonSerialization) Serialization() BaseSerializationImpl {
	return BaseSerializationImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = JsonSerialization{}

func (s JsonSerialization) MarshalJSON() ([]byte, error) {
	type wrapper JsonSerialization
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling JsonSerialization: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling JsonSerialization: %+v", err)
	}

	decoded["type"] = "Json"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling JsonSerialization: %+v", err)
	}

	return encoded, nil
}
