package datasets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExcelDatasetTypeProperties struct {
	Compression      *DatasetCompression `json:"compression,omitempty"`
	FirstRowAsHeader *bool               `json:"firstRowAsHeader,omitempty"`
	Location         DatasetLocation     `json:"location"`
	NullValue        *string             `json:"nullValue,omitempty"`
	Range            *string             `json:"range,omitempty"`
	SheetIndex       *int64              `json:"sheetIndex,omitempty"`
	SheetName        *string             `json:"sheetName,omitempty"`
}

var _ json.Unmarshaler = &ExcelDatasetTypeProperties{}

func (s *ExcelDatasetTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Compression      *DatasetCompression `json:"compression,omitempty"`
		FirstRowAsHeader *bool               `json:"firstRowAsHeader,omitempty"`
		NullValue        *string             `json:"nullValue,omitempty"`
		Range            *string             `json:"range,omitempty"`
		SheetIndex       *int64              `json:"sheetIndex,omitempty"`
		SheetName        *string             `json:"sheetName,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Compression = decoded.Compression
	s.FirstRowAsHeader = decoded.FirstRowAsHeader
	s.NullValue = decoded.NullValue
	s.Range = decoded.Range
	s.SheetIndex = decoded.SheetIndex
	s.SheetName = decoded.SheetName

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ExcelDatasetTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["location"]; ok {
		impl, err := UnmarshalDatasetLocationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Location' for 'ExcelDatasetTypeProperties': %+v", err)
		}
		s.Location = impl
	}

	return nil
}
