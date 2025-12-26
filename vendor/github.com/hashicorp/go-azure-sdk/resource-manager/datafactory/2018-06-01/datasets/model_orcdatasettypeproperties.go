package datasets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OrcDatasetTypeProperties struct {
	Location            DatasetLocation `json:"location"`
	OrcCompressionCodec *interface{}    `json:"orcCompressionCodec,omitempty"`
}

var _ json.Unmarshaler = &OrcDatasetTypeProperties{}

func (s *OrcDatasetTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		OrcCompressionCodec *interface{} `json:"orcCompressionCodec,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.OrcCompressionCodec = decoded.OrcCompressionCodec

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling OrcDatasetTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["location"]; ok {
		impl, err := UnmarshalDatasetLocationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Location' for 'OrcDatasetTypeProperties': %+v", err)
		}
		s.Location = impl
	}

	return nil
}
