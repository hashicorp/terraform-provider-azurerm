package datasets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DatasetStorageFormat = JsonFormat{}

type JsonFormat struct {
	EncodingName       *string      `json:"encodingName,omitempty"`
	FilePattern        *interface{} `json:"filePattern,omitempty"`
	JsonNodeReference  *string      `json:"jsonNodeReference,omitempty"`
	JsonPathDefinition *interface{} `json:"jsonPathDefinition,omitempty"`
	NestingSeparator   *string      `json:"nestingSeparator,omitempty"`

	// Fields inherited from DatasetStorageFormat

	Deserializer *string `json:"deserializer,omitempty"`
	Serializer   *string `json:"serializer,omitempty"`
	Type         string  `json:"type"`
}

func (s JsonFormat) DatasetStorageFormat() BaseDatasetStorageFormatImpl {
	return BaseDatasetStorageFormatImpl{
		Deserializer: s.Deserializer,
		Serializer:   s.Serializer,
		Type:         s.Type,
	}
}

var _ json.Marshaler = JsonFormat{}

func (s JsonFormat) MarshalJSON() ([]byte, error) {
	type wrapper JsonFormat
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling JsonFormat: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling JsonFormat: %+v", err)
	}

	decoded["type"] = "JsonFormat"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling JsonFormat: %+v", err)
	}

	return encoded, nil
}
