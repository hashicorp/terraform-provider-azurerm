package datasets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DatasetStorageFormat = AvroFormat{}

type AvroFormat struct {

	// Fields inherited from DatasetStorageFormat

	Deserializer *string `json:"deserializer,omitempty"`
	Serializer   *string `json:"serializer,omitempty"`
	Type         string  `json:"type"`
}

func (s AvroFormat) DatasetStorageFormat() BaseDatasetStorageFormatImpl {
	return BaseDatasetStorageFormatImpl{
		Deserializer: s.Deserializer,
		Serializer:   s.Serializer,
		Type:         s.Type,
	}
}

var _ json.Marshaler = AvroFormat{}

func (s AvroFormat) MarshalJSON() ([]byte, error) {
	type wrapper AvroFormat
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AvroFormat: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AvroFormat: %+v", err)
	}

	decoded["type"] = "AvroFormat"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AvroFormat: %+v", err)
	}

	return encoded, nil
}
