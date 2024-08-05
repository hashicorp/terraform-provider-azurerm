package alertsmanagements

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertsMetaData struct {
	Properties AlertsMetaDataProperties `json:"properties"`
}

var _ json.Unmarshaler = &AlertsMetaData{}

func (s *AlertsMetaData) UnmarshalJSON(bytes []byte) error {

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AlertsMetaData into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["properties"]; ok {
		impl, err := unmarshalAlertsMetaDataPropertiesImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Properties' for 'AlertsMetaData': %+v", err)
		}
		s.Properties = impl
	}
	return nil
}
