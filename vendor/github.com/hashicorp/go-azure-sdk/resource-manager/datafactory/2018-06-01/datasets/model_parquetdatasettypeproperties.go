package datasets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ParquetDatasetTypeProperties struct {
	CompressionCodec *interface{}    `json:"compressionCodec,omitempty"`
	Location         DatasetLocation `json:"location"`
}

var _ json.Unmarshaler = &ParquetDatasetTypeProperties{}

func (s *ParquetDatasetTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		CompressionCodec *interface{} `json:"compressionCodec,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.CompressionCodec = decoded.CompressionCodec

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ParquetDatasetTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["location"]; ok {
		impl, err := UnmarshalDatasetLocationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Location' for 'ParquetDatasetTypeProperties': %+v", err)
		}
		s.Location = impl
	}

	return nil
}
