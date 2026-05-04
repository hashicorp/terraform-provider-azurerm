package streamingjobs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Serialization = CsvSerialization{}

type CsvSerialization struct {
	Properties *CsvSerializationProperties `json:"properties,omitempty"`

	// Fields inherited from Serialization

	Type EventSerializationType `json:"type"`
}

func (s CsvSerialization) Serialization() BaseSerializationImpl {
	return BaseSerializationImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = CsvSerialization{}

func (s CsvSerialization) MarshalJSON() ([]byte, error) {
	type wrapper CsvSerialization
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CsvSerialization: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CsvSerialization: %+v", err)
	}

	decoded["type"] = "Csv"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CsvSerialization: %+v", err)
	}

	return encoded, nil
}
