package streamingjobs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OutputProperties struct {
	Datasource    OutputDataSource `json:"datasource"`
	Diagnostics   *Diagnostics     `json:"diagnostics,omitempty"`
	Etag          *string          `json:"etag,omitempty"`
	Serialization Serialization    `json:"serialization"`
	SizeWindow    *int64           `json:"sizeWindow,omitempty"`
	TimeWindow    *string          `json:"timeWindow,omitempty"`
}

var _ json.Unmarshaler = &OutputProperties{}

func (s *OutputProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Diagnostics *Diagnostics `json:"diagnostics,omitempty"`
		Etag        *string      `json:"etag,omitempty"`
		SizeWindow  *int64       `json:"sizeWindow,omitempty"`
		TimeWindow  *string      `json:"timeWindow,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Diagnostics = decoded.Diagnostics
	s.Etag = decoded.Etag
	s.SizeWindow = decoded.SizeWindow
	s.TimeWindow = decoded.TimeWindow

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling OutputProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["datasource"]; ok {
		impl, err := UnmarshalOutputDataSourceImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Datasource' for 'OutputProperties': %+v", err)
		}
		s.Datasource = impl
	}

	if v, ok := temp["serialization"]; ok {
		impl, err := UnmarshalSerializationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Serialization' for 'OutputProperties': %+v", err)
		}
		s.Serialization = impl
	}

	return nil
}
