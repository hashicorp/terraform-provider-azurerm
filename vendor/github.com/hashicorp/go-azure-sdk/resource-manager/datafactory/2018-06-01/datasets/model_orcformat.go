package datasets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DatasetStorageFormat = OrcFormat{}

type OrcFormat struct {

	// Fields inherited from DatasetStorageFormat

	Deserializer *interface{} `json:"deserializer,omitempty"`
	Serializer   *interface{} `json:"serializer,omitempty"`
	Type         string       `json:"type"`
}

func (s OrcFormat) DatasetStorageFormat() BaseDatasetStorageFormatImpl {
	return BaseDatasetStorageFormatImpl{
		Deserializer: s.Deserializer,
		Serializer:   s.Serializer,
		Type:         s.Type,
	}
}

var _ json.Marshaler = OrcFormat{}

func (s OrcFormat) MarshalJSON() ([]byte, error) {
	type wrapper OrcFormat
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling OrcFormat: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling OrcFormat: %+v", err)
	}

	decoded["type"] = "OrcFormat"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling OrcFormat: %+v", err)
	}

	return encoded, nil
}
