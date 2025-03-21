package datasets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DelimitedTextDatasetTypeProperties struct {
	ColumnDelimiter  *string         `json:"columnDelimiter,omitempty"`
	CompressionCodec *string         `json:"compressionCodec,omitempty"`
	CompressionLevel *string         `json:"compressionLevel,omitempty"`
	EncodingName     *string         `json:"encodingName,omitempty"`
	EscapeChar       *string         `json:"escapeChar,omitempty"`
	FirstRowAsHeader *bool           `json:"firstRowAsHeader,omitempty"`
	Location         DatasetLocation `json:"location"`
	NullValue        *string         `json:"nullValue,omitempty"`
	QuoteChar        *string         `json:"quoteChar,omitempty"`
	RowDelimiter     *string         `json:"rowDelimiter,omitempty"`
}

var _ json.Unmarshaler = &DelimitedTextDatasetTypeProperties{}

func (s *DelimitedTextDatasetTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ColumnDelimiter  *string `json:"columnDelimiter,omitempty"`
		CompressionCodec *string `json:"compressionCodec,omitempty"`
		CompressionLevel *string `json:"compressionLevel,omitempty"`
		EncodingName     *string `json:"encodingName,omitempty"`
		EscapeChar       *string `json:"escapeChar,omitempty"`
		FirstRowAsHeader *bool   `json:"firstRowAsHeader,omitempty"`
		NullValue        *string `json:"nullValue,omitempty"`
		QuoteChar        *string `json:"quoteChar,omitempty"`
		RowDelimiter     *string `json:"rowDelimiter,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ColumnDelimiter = decoded.ColumnDelimiter
	s.CompressionCodec = decoded.CompressionCodec
	s.CompressionLevel = decoded.CompressionLevel
	s.EncodingName = decoded.EncodingName
	s.EscapeChar = decoded.EscapeChar
	s.FirstRowAsHeader = decoded.FirstRowAsHeader
	s.NullValue = decoded.NullValue
	s.QuoteChar = decoded.QuoteChar
	s.RowDelimiter = decoded.RowDelimiter

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling DelimitedTextDatasetTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["location"]; ok {
		impl, err := UnmarshalDatasetLocationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Location' for 'DelimitedTextDatasetTypeProperties': %+v", err)
		}
		s.Location = impl
	}

	return nil
}
