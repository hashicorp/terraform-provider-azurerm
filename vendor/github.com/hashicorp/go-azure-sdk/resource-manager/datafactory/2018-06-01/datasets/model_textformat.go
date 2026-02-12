package datasets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DatasetStorageFormat = TextFormat{}

type TextFormat struct {
	ColumnDelimiter  *interface{} `json:"columnDelimiter,omitempty"`
	EncodingName     *interface{} `json:"encodingName,omitempty"`
	EscapeChar       *interface{} `json:"escapeChar,omitempty"`
	FirstRowAsHeader *bool        `json:"firstRowAsHeader,omitempty"`
	NullValue        *interface{} `json:"nullValue,omitempty"`
	QuoteChar        *interface{} `json:"quoteChar,omitempty"`
	RowDelimiter     *interface{} `json:"rowDelimiter,omitempty"`
	SkipLineCount    *int64       `json:"skipLineCount,omitempty"`
	TreatEmptyAsNull *bool        `json:"treatEmptyAsNull,omitempty"`

	// Fields inherited from DatasetStorageFormat

	Deserializer *interface{} `json:"deserializer,omitempty"`
	Serializer   *interface{} `json:"serializer,omitempty"`
	Type         string       `json:"type"`
}

func (s TextFormat) DatasetStorageFormat() BaseDatasetStorageFormatImpl {
	return BaseDatasetStorageFormatImpl{
		Deserializer: s.Deserializer,
		Serializer:   s.Serializer,
		Type:         s.Type,
	}
}

var _ json.Marshaler = TextFormat{}

func (s TextFormat) MarshalJSON() ([]byte, error) {
	type wrapper TextFormat
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling TextFormat: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling TextFormat: %+v", err)
	}

	decoded["type"] = "TextFormat"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling TextFormat: %+v", err)
	}

	return encoded, nil
}
