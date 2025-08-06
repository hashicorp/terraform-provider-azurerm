package outputs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Serialization = ParquetSerialization{}

type ParquetSerialization struct {
	Properties *interface{} `json:"properties,omitempty"`

	// Fields inherited from Serialization

	Type EventSerializationType `json:"type"`
}

func (s ParquetSerialization) Serialization() BaseSerializationImpl {
	return BaseSerializationImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = ParquetSerialization{}

func (s ParquetSerialization) MarshalJSON() ([]byte, error) {
	type wrapper ParquetSerialization
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ParquetSerialization: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ParquetSerialization: %+v", err)
	}

	decoded["type"] = "Parquet"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ParquetSerialization: %+v", err)
	}

	return encoded, nil
}
