package datasets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IcebergDatasetTypeProperties struct {
	Location DatasetLocation `json:"location"`
}

var _ json.Unmarshaler = &IcebergDatasetTypeProperties{}

func (s *IcebergDatasetTypeProperties) UnmarshalJSON(bytes []byte) error {

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling IcebergDatasetTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["location"]; ok {
		impl, err := UnmarshalDatasetLocationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Location' for 'IcebergDatasetTypeProperties': %+v", err)
		}
		s.Location = impl
	}

	return nil
}
