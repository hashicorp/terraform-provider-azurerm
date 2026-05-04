package streamingjobs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Serialization = CustomClrSerialization{}

type CustomClrSerialization struct {
	Properties *CustomClrSerializationProperties `json:"properties,omitempty"`

	// Fields inherited from Serialization

	Type EventSerializationType `json:"type"`
}

func (s CustomClrSerialization) Serialization() BaseSerializationImpl {
	return BaseSerializationImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = CustomClrSerialization{}

func (s CustomClrSerialization) MarshalJSON() ([]byte, error) {
	type wrapper CustomClrSerialization
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CustomClrSerialization: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CustomClrSerialization: %+v", err)
	}

	decoded["type"] = "CustomClr"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CustomClrSerialization: %+v", err)
	}

	return encoded, nil
}
