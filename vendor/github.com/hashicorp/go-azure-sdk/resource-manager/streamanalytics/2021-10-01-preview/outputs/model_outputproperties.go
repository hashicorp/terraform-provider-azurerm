package outputs

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OutputProperties struct {
	Datasource                OutputDataSource            `json:"datasource"`
	Diagnostics               *Diagnostics                `json:"diagnostics,omitempty"`
	Etag                      *string                     `json:"etag,omitempty"`
	LastOutputEventTimestamps *[]LastOutputEventTimestamp `json:"lastOutputEventTimestamps,omitempty"`
	Serialization             Serialization               `json:"serialization"`
	SizeWindow                *int64                      `json:"sizeWindow,omitempty"`
	TimeWindow                *string                     `json:"timeWindow,omitempty"`
	WatermarkSettings         *OutputWatermarkProperties  `json:"watermarkSettings,omitempty"`
}

var _ json.Unmarshaler = &OutputProperties{}

func (s *OutputProperties) UnmarshalJSON(bytes []byte) error {
	type alias OutputProperties
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into OutputProperties: %+v", err)
	}

	s.Diagnostics = decoded.Diagnostics
	s.Etag = decoded.Etag
	s.LastOutputEventTimestamps = decoded.LastOutputEventTimestamps
	s.SizeWindow = decoded.SizeWindow
	s.TimeWindow = decoded.TimeWindow
	s.WatermarkSettings = decoded.WatermarkSettings

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling OutputProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["datasource"]; ok {
		impl, err := unmarshalOutputDataSourceImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Datasource' for 'OutputProperties': %+v", err)
		}
		s.Datasource = impl
	}

	if v, ok := temp["serialization"]; ok {
		impl, err := unmarshalSerializationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Serialization' for 'OutputProperties': %+v", err)
		}
		s.Serialization = impl
	}
	return nil
}
