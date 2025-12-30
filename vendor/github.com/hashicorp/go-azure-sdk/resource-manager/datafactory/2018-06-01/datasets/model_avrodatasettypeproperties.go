package datasets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvroDatasetTypeProperties struct {
	AvroCompressionCodec *interface{}    `json:"avroCompressionCodec,omitempty"`
	AvroCompressionLevel *int64          `json:"avroCompressionLevel,omitempty"`
	Location             DatasetLocation `json:"location"`
}

var _ json.Unmarshaler = &AvroDatasetTypeProperties{}

func (s *AvroDatasetTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AvroCompressionCodec *interface{} `json:"avroCompressionCodec,omitempty"`
		AvroCompressionLevel *int64       `json:"avroCompressionLevel,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AvroCompressionCodec = decoded.AvroCompressionCodec
	s.AvroCompressionLevel = decoded.AvroCompressionLevel

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AvroDatasetTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["location"]; ok {
		impl, err := UnmarshalDatasetLocationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Location' for 'AvroDatasetTypeProperties': %+v", err)
		}
		s.Location = impl
	}

	return nil
}
