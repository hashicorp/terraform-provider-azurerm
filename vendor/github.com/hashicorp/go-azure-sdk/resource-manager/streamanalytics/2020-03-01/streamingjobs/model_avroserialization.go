package streamingjobs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Serialization = AvroSerialization{}

type AvroSerialization struct {
	Properties *interface{} `json:"properties,omitempty"`

	// Fields inherited from Serialization

	Type EventSerializationType `json:"type"`
}

func (s AvroSerialization) Serialization() BaseSerializationImpl {
	return BaseSerializationImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = AvroSerialization{}

func (s AvroSerialization) MarshalJSON() ([]byte, error) {
	type wrapper AvroSerialization
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AvroSerialization: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AvroSerialization: %+v", err)
	}

	decoded["type"] = "Avro"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AvroSerialization: %+v", err)
	}

	return encoded, nil
}
