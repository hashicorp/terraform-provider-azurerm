package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetMetadataActivityTypeProperties struct {
	Dataset        DatasetReference   `json:"dataset"`
	FieldList      *[]interface{}     `json:"fieldList,omitempty"`
	FormatSettings FormatReadSettings `json:"formatSettings"`
	StoreSettings  StoreReadSettings  `json:"storeSettings"`
}

var _ json.Unmarshaler = &GetMetadataActivityTypeProperties{}

func (s *GetMetadataActivityTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Dataset   DatasetReference `json:"dataset"`
		FieldList *[]interface{}   `json:"fieldList,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Dataset = decoded.Dataset
	s.FieldList = decoded.FieldList

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling GetMetadataActivityTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["formatSettings"]; ok {
		impl, err := UnmarshalFormatReadSettingsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'FormatSettings' for 'GetMetadataActivityTypeProperties': %+v", err)
		}
		s.FormatSettings = impl
	}

	if v, ok := temp["storeSettings"]; ok {
		impl, err := UnmarshalStoreReadSettingsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'StoreSettings' for 'GetMetadataActivityTypeProperties': %+v", err)
		}
		s.StoreSettings = impl
	}

	return nil
}
